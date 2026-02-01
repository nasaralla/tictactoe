package play

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Player struct {
	Connection net.Conn
	Reader     *bufio.Reader
	Writer     *bufio.Writer
	Mark       byte
	Series     []int
}

type Game struct {
	Board   [9]byte
	Players [2]*Player
	Turn    int
}

func getInputValue(input string) int {
	s := strings.Split(input, " ")
	i, err := strconv.Atoi(s[1])
	if err != nil {
		// handle error
		panic(err)
	}
	return i
}

func isValidInput(input int, game *Game) bool {
	if input < 0 || input > 8 || string(game.Board[input]) != "." {
		return false
	}
	return true
}

func RunGame(game *Game) {
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Player 1 Loaded:", game.Players[0].Mark)
	fmt.Println("Player 2 Loaded:", game.Players[1].Mark)

	for _, p := range game.Players {
		p.Writer.WriteString(fmt.Sprintf("WELCOME %c\n", p.Mark))
		p.Writer.Flush()
	}

	for _, p := range game.Players {
		p.Writer.WriteString("START\n")
		p.Writer.Flush()
	}

	for {
		current := game.Players[game.Turn]
		other := game.Players[1-game.Turn]

		sendBoard(game)

		current.Writer.WriteString("YOUR_TURN\n")
		current.Writer.Flush()

		line, err := current.Reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "MOVE ") {
			current.Writer.WriteString("ERROR Invalid command\n")
			current.Writer.Flush()
			continue
		}

		pos := getInputValue(line)

		if isValidInput(pos, game) {
			current.Writer.WriteString("ERROR Invalid move\n")
			current.Writer.Flush()
			continue
		}

		game.Board[pos] = current.Mark
		other.Writer.WriteString(fmt.Sprintf("OPPONENT_MOVED %d\n", pos))
		other.Writer.Flush()

		if winner(game.Board, current.Mark) {
			current.Writer.WriteString("WIN\n")
			other.Writer.WriteString("LOSE\n")
			current.Writer.Flush()
			other.Writer.Flush()
			return
		}

		if full(game.Board) {
			for _, p := range game.Players {
				p.Writer.WriteString("DRAW\n")
				p.Writer.Flush()
			}
			return
		}

		game.Turn = 1 - game.Turn
	}
}

func sendBoard(game *Game) {
	Board := string(game.Board[:])
	for _, p := range game.Players {
		p.Writer.WriteString("BOARD " + Board + "\n")
		p.Writer.Flush()
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
