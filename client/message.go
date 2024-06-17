package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type      string `json:"type"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Room      string `json:"room,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Password  string `json:"password,omitempty"`
}

func readMessages(conn *websocket.Conn) {
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			return
		}
		messages = append(messages, msg)

		if msg.Recipient != "" && msg.Recipient == username {
			fmt.Printf("[DM from %s] %s\n", msg.Username, msg.Content)
		} else if msg.Room == "" || msg.Room == room {
			fmt.Printf("[%s] %s: %s\n", msg.Room, msg.Username, msg.Content)
		}
	}
}

func sendMessage(conn *websocket.Conn, content string) {
	msg := Message{
		Type:     "chat",
		Username: username,
		Content:  content,
		Room:     room,
	}
	err := conn.WriteJSON(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func sendDirectMessage(conn *websocket.Conn, recipient, content string) {
	msg := Message{
		Type:      "dm",
		Username:  username,
		Recipient: recipient,
		Content:   content,
	}
	err := conn.WriteJSON(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	fmt.Printf("[DM to %s]: %s\n", recipient, content)
}

func saveMessages() {
	if _, err := os.Stat("message_logs"); os.IsNotExist(err) {
		err := os.Mkdir("message_logs", 0755)
		if err != nil {
			log.Fatalf("failed to create message_logs directory: %v", err)
		}
	}
	fileName := fmt.Sprintf("message_logs/%s_%s.json", username, time.Now().Format("20060102_150405"))
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("failed to create log file: %v", err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal messages: %v", err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("failed to write messages to log file: %v", err)
	}

	fmt.Printf("Messages saved to %s\n", fileName)
}
