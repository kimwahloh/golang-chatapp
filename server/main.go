package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[string]*Client)
	rooms     = make(map[string]map[string]*Client)
	broadcast = make(chan Message)
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	roomFile     = "room.json"
	usersFile    = "users.json"
	users        = make(map[string]string)
	usersMutex   sync.Mutex
	clientsMutex sync.Mutex
	roomsMutex   sync.Mutex
)

func main() {
	flag.Parse()
	loadRoomsFromJSON()
	loadUsersFromJSON()

	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Println("HTTP server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
