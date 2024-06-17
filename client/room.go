package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func joinRoom(conn *websocket.Conn, roomName string) {
	msg := Message{
		Type:     "join",
		Username: username,
		Room:     roomName,
	}
	err := conn.WriteJSON(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}
	room = roomName
}

func leaveRoom(conn *websocket.Conn) {
	msg := Message{
		Type:     "leave",
		Username: username,
		Room:     room,
	}
	err := conn.WriteJSON(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}
	room = ""

	joinRoom(conn, "public")
}

func createRoom(conn *websocket.Conn, roomName string) {
	msg := Message{
		Type:     "create",
		Username: username,
		Room:     roomName,
	}
	err := conn.WriteJSON(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}
	fmt.Printf("Room %s created.\n", roomName)
	joinRoom(conn, roomName)
}
