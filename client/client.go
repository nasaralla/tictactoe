package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("ERR:", err)
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	stdin := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("ERR 2:", err)
			return
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")

		fmt.Println("Line:", line)

		switch parts[0] {
		case "WELCOME":
			fmt.Println("You are player", parts[1])

		case "START":
			fmt.Println("Game started")

		case "BOARD":
			printBoard(parts[1])

		case "YOUR_TURN":
			fmt.Print("Your move (0-8): ")
			input, _ := stdin.ReadString('\n')
			input = strings.TrimSpace(input)
			writer.WriteString("MOVE " + input + "\n")
			writer.Flush()

		case "OPPONENT_MOVED":
			fmt.Println("Opponent moved at", parts[1])

		case "WIN":
			fmt.Println("You win!")
			return

		case "LOSE":
			fmt.Println("You lose!")
			return

		case "DRAW":
			fmt.Println("Draw!")
			return

		case "ERROR":
			fmt.Println("Error:", strings.Join(parts[1:], " "))
		}
	}
}

func printBoard(b string) {
	fmt.Println()
	for i := 0; i < 9; i += 3 {
		fmt.Printf(" %c | %c | %c \n", b[i], b[i+1], b[i+2])
		if i < 6 {
			fmt.Println("---+---+---")
		}
	}
	fmt.Println()
}
