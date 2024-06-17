package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	serverAddr = flag.String("addr", "localhost:8000", "http service address")
	username   string
	room       string
	messages   []Message
)

func main() {
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Choose action: (1) Register (2) Login: ")
		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(action)

		if action == "1" {
			register(reader)
			break
		} else if action == "2" {
			login(reader)
			break
		} else {
			fmt.Println("Invalid choice. Please choose 1 or 2.")
		}
	}

	connectToServer()
}
