package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func connectToServer() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+*serverAddr+"/ws?username="+username, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	joinRoom(conn, "public")

	go readMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		switch {
		case strings.HasPrefix(input, "/create"):
			roomName := strings.TrimSpace(strings.TrimPrefix(input, "/create"))
			createRoom(conn, roomName)
		case strings.HasPrefix(input, "/join"):
			roomName := strings.TrimSpace(strings.TrimPrefix(input, "/join"))
			joinRoom(conn, roomName)
		case input == "/leave":
			leaveRoom(conn)
		case input == "/save":
			saveMessages()
		case strings.HasPrefix(input, "/dm"):
			parts := strings.SplitN(input, " ", 3)
			if len(parts) < 3 {
				fmt.Println("Usage: /dm <recipient> <message>")
				continue
			}
			recipient := parts[1]
			message := parts[2]
			sendDirectMessage(conn, recipient, message)
		default:
			sendMessage(conn, input)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading standard input: %v", err)
	}
}
