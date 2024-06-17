package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

func register(reader *bufio.Reader) {
	fmt.Print("Enter your username: ")
	usernameInput, _ := reader.ReadString('\n')
	username = strings.TrimSpace(usernameInput)

	fmt.Print("Enter your password: ")
	passwordInput, _ := reader.ReadString('\n')
	password := strings.TrimSpace(passwordInput)

	conn, _, err := websocket.DefaultDialer.Dial("ws://"+*serverAddr+"/ws?username="+username, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	msg := Message{
		Type:     "register",
		Username: username,
		Password: password,
	}
	err = conn.WriteJSON(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	var resp Message
	err = conn.ReadJSON(&resp)
	if err != nil {
		log.Println("read:", err)
		return
	}

	if strings.Contains(resp.Content, "successful") {
		fmt.Println("Registration successful. Please login.")
		login(reader)
	} else {
		fmt.Println(resp.Content)
		register(reader)
	}
}

func login(reader *bufio.Reader) {
	fmt.Print("Enter your username: ")
	usernameInput, _ := reader.ReadString('\n')
	username = strings.TrimSpace(usernameInput)

	fmt.Print("Enter your password: ")
	passwordInput, _ := reader.ReadString('\n')
	password := strings.TrimSpace(passwordInput)

	conn, _, err := websocket.DefaultDialer.Dial("ws://"+*serverAddr+"/ws?username="+username, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	msg := Message{
		Type:     "login",
		Username: username,
		Password: password,
	}
	err = conn.WriteJSON(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	var resp Message
	err = conn.ReadJSON(&resp)
	if err != nil {
		log.Println("read:", err)
		return
	}

	if strings.Contains(resp.Content, "successful") {
		fmt.Println("Login successful.")

	} else {
		fmt.Println(resp.Content)
		login(reader)
	}
}
