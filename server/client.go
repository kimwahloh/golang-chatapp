package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Username string
	Conn     *websocket.Conn
	Rooms    map[string]bool
}

func closeConnection(conn *websocket.Conn) {
	err := conn.Close()
	if err != nil {
		log.Printf("Error closing connection: %v", err)
	} else {
		log.Println("Connection closed successfully")
	}

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for username, client := range clients {
		if client.Conn == conn {
			for room := range client.Rooms {
				delete(rooms[room], username)
				notifyRoomUpdate(room)
			}
			delete(clients, username)
			break
		}
	}
}
