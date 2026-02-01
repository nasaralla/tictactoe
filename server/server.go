package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Player struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	mark   byte
}

type Game struct {
	board   [9]byte
	players [2]*Player
	turn    int
}

var waitingPlayers = make(chan *Player)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Server listening on :9000")

	go matchMaker()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	player := &Player{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
	waitingPlayers <- player
}

func matchMaker() {
	for {
		p1 := <-waitingPlayers
		p2 := <-waitingPlayers

		p1.mark = 'X'
		p2.mark = 'O'

		game := &Game{
			players: [2]*Player{p1, p2},
			turn:    0,
		}

		for i := range game.board {
			game.board[i] = '.'
		}

		go runGame(game)
	}
}

func runGame(g *Game) {
	var wg sync.WaitGroup
	wg.Add(2)

	for _, p := range g.players {
		p.writer.WriteString(fmt.Sprintf("WELCOME %c\n", p.mark))
		p.writer.Flush()
	}

	for _, p := range g.players {
		p.writer.WriteString("START\n")
		p.writer.Flush()
	}

	for {
		current := g.players[g.turn]
		other := g.players[1-g.turn]

		sendBoard(g)

		current.writer.WriteString("YOUR_TURN\n")
		current.writer.Flush()

		line, err := current.reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "MOVE ") {
			current.writer.WriteString("ERROR Invalid command\n")
			current.writer.Flush()
			continue
		}

		var pos int
		fmt.Sscanf(line, "MOVE %d", &pos)

		if pos < 0 || pos > 8 || g.board[pos] != '.' {
			current.writer.WriteString("ERROR Invalid move\n")
			current.writer.Flush()
			continue
		}

		g.board[pos] = current.mark
		other.writer.WriteString(fmt.Sprintf("OPPONENT_MOVED %d\n", pos))
		other.writer.Flush()

		if winner(g.board, current.mark) {
			current.writer.WriteString("WIN\n")
			other.writer.WriteString("LOSE\n")
			current.writer.Flush()
			other.writer.Flush()
			return
		}

		if full(g.board) {
			for _, p := range g.players {
				p.writer.WriteString("DRAW\n")
				p.writer.Flush()
			}
			return
		}

		g.turn = 1 - g.turn
	}
}

func sendBoard(g *Game) {
	board := string(g.board[:])
	for _, p := range g.players {
		p.writer.WriteString("BOARD " + board + "\n")
		p.writer.Flush()
	}
}

func winner(b [9]byte, m byte) bool {
	wins := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}
	for _, w := range wins {
		if b[w[0]] == m && b[w[1]] == m && b[w[2]] == m {
			return true
		}
	}
	return false
}

func full(b [9]byte) bool {
	for _, c := range b {
		if c == '.' {
			return false
		}
	}
	return true
}
