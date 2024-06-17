package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("Unauthorized connection attempt")
		return
	}

	client := &Client{
		Username: username,
		Conn:     ws,
		Rooms:    make(map[string]bool),
	}
	clientsMutex.Lock()
	clients[username] = client
	clientsMutex.Unlock()
	log.Println("New client connected:", username)

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			clientsMutex.Lock()
			delete(clients, username)
			clientsMutex.Unlock()
			closeConnection(ws)
			return
		}

		switch msg.Type {
		case "register":
			handleRegister(ws, msg)
		case "login":
			handleLogin(ws, msg)
		default:
			client, ok := clients[msg.Username]
			if !ok {
				ws.WriteJSON(Message{
					Type:    "info",
					Content: "You must be logged in to perform this action.",
				})
				continue
			}

			if strings.HasPrefix(msg.Content, "/") {
				handleCommand(msg, client)
			} else {
				switch msg.Type {
				case "join":
					handleJoinRoom(msg, client)
				case "leave":
					handleLeaveRoom(msg, client)
				case "create":
					handleCreateRoom(msg, client)
				case "dm":
					handleDirectMessage(msg, client)
				default:
					broadcast <- msg
				}
			}
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		log.Printf("Received message: %s: %s", msg.Username, msg.Content)

		if msg.Room != "" {

			roomClients, ok := rooms[msg.Room]
			if !ok {
				log.Printf("Room %s not found", msg.Room)
				continue
			}
			for _, client := range roomClients {
				if client.Rooms[msg.Room] {
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						log.Printf("error sending message to %s: %v", client.Username, err)
						closeConnection(client.Conn)
						delete(clients, client.Username)
					}
				}
			}
		} else {

			for _, client := range clients {
				if len(client.Rooms) == 0 {
					shouldBroadcast := true
					for room := range rooms {
						if _, ok := rooms[room][client.Username]; ok {
							shouldBroadcast = false
							break
						}
					}
					if shouldBroadcast {
						err := client.Conn.WriteJSON(msg)
						if err != nil {
							log.Printf("error broadcasting message to %s: %v", client.Username, err)
							closeConnection(client.Conn)
							delete(clients, client.Username)
						}
					}
				}
			}
		}
	}
}

func handleRegister(ws *websocket.Conn, msg Message) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	if _, exists := users[msg.Username]; exists {
		ws.WriteJSON(Message{
			Type:    "info",
			Content: "Username already exists.",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(msg.Password), bcrypt.DefaultCost)
	if err != nil {
		ws.WriteJSON(Message{
			Type:    "info",
			Content: "Error registering user.",
		})
		return
	}

	users[msg.Username] = string(hashedPassword)
	saveUsersToJSON(users)

	ws.WriteJSON(Message{
		Type:    "info",
		Content: "Registration successful.",
	})
}

func handleLogin(ws *websocket.Conn, msg Message) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	storedPassword, exists := users[msg.Username]
	if !exists || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(msg.Password)) != nil {
		ws.WriteJSON(Message{
			Type:    "info",
			Content: "Invalid username or password.",
		})
		return
	}

	client := &Client{
		Username: msg.Username,
		Conn:     ws,
		Rooms:    make(map[string]bool),
	}

	clientsMutex.Lock()
	clients[msg.Username] = client
	clientsMutex.Unlock()

	ws.WriteJSON(Message{
		Type:    "info",
		Content: "Login successful.",
	})

	handleJoinRoom(Message{
		Type:     "join",
		Username: msg.Username,
		Room:     "public",
	}, client)
}

func handleCommand(msg Message, client *Client) {
	command := strings.TrimPrefix(msg.Content, "/")
	if command == "" {
		client.Conn.WriteJSON(Message{
			Type:    "info",
			Content: "Available commands: /list, /join {room_name}, /leave {room_name}, /create {room_name}, /dm {recipient} {message}",
		})
		return
	}
	parts := strings.SplitN(command, " ", 2)
	switch parts[0] {
	case "list":
		handleListRooms(client)
	case "join":
		if len(parts) < 2 {
			client.Conn.WriteJSON(Message{
				Type:    "info",
				Content: "Usage: /join room_name",
			})
			return
		}
		handleJoinRoom(Message{
			Type:     "join",
			Username: client.Username,
			Room:     parts[1],
		}, client)
	case "leave":
		if len(parts) < 2 {
			client.Conn.WriteJSON(Message{
				Type:    "info",
				Content: "Usage: /leave room_name",
			})
			return
		}
		handleLeaveRoom(Message{
			Type:     "leave",
			Username: client.Username,
			Room:     parts[1],
		}, client)
	case "create":
		if len(parts) < 2 {
			client.Conn.WriteJSON(Message{
				Type:    "info",
				Content: "Usage: /create room_name",
			})
			return
		}
		handleCreateRoom(Message{
			Type:     "create",
			Username: client.Username,
			Room:     parts[1],
		}, client)
	case "dm":
		if len(parts) < 2 {
			client.Conn.WriteJSON(Message{
				Type:    "info",
				Content: "Usage: /dm recipient message",
			})
			return
		}
		dmParts := strings.SplitN(parts[1], " ", 2)
		if len(dmParts) < 2 {
			client.Conn.WriteJSON(Message{
				Type:    "info",
				Content: "Usage: /dm recipient message",
			})
			return
		}
		handleDirectMessage(Message{
			Type:      "dm",
			Username:  client.Username,
			Recipient: dmParts[0],
			Content:   dmParts[1],
		}, client)
	default:
		client.Conn.WriteJSON(Message{
			Type:    "info",
			Content: "Unknown command",
		})
	}
}

func handleListRooms(client *Client) {
	roomNames := make([]string, 0, len(rooms))
	for roomName := range rooms {
		roomNames = append(roomNames, roomName)
	}
	client.Conn.WriteJSON(Message{
		Type:    "info",
		Content: "Rooms: " + strings.Join(roomNames, ", "),
	})
}

func handleJoinRoom(msg Message, client *Client) {
	roomName := msg.Room

	if _, ok := rooms[roomName]; !ok {
		client.Conn.WriteJSON(Message{
			Type:    "error",
			Content: fmt.Sprintf("Room %s does not exist. Please create the room first.", roomName),
		})
		log.Printf("Room %s does not exist. User %s needs to create the room first.", roomName, client.Username)
		return
	}

	rooms[roomName][client.Username] = client
	client.Rooms[roomName] = true

	notifyRoomUpdate(roomName)

	saveRoomsToJSON(generateRoomsJSON())

	log.Printf("User %s joined room %s", client.Username, roomName)
}

func handleLeaveRoom(msg Message, client *Client) {
	roomName := msg.Room
	if _, ok := rooms[roomName]; !ok {
		log.Printf("Room %s not found", roomName)
		return
	}
	delete(rooms[roomName], client.Username)
	delete(client.Rooms, roomName)

	notifyRoomUpdate(roomName)

	saveRoomsToJSON(generateRoomsJSON())

	if _, ok := rooms["public"]; !ok {
		rooms["public"] = make(map[string]*Client)
	}
	rooms["public"][client.Username] = client
	client.Rooms["public"] = true

	notifyRoomUpdate("public")

	client.Conn.WriteJSON(Message{
		Type:    "info",
		Content: fmt.Sprintf("You have left the room %s and joined the public room", roomName),
	})

	log.Printf("User %s left room %s and joined public room", client.Username, roomName)
}

func handleCreateRoom(msg Message, client *Client) {
	roomName := msg.Room
	if _, ok := rooms[roomName]; ok {
		log.Printf("Room %s already exists", roomName)
		client.Conn.WriteJSON(Message{
			Type:    "info",
			Content: fmt.Sprintf("Room %s already exists", roomName),
		})
		return
	}
	rooms[roomName] = make(map[string]*Client)
	rooms[roomName][client.Username] = client
	client.Rooms[roomName] = true

	saveRoomsToJSON(generateRoomsJSON())
}

func handleDirectMessage(msg Message, client *Client) {
	recipientClient, ok := clients[msg.Recipient]
	if !ok {
		client.Conn.WriteJSON(Message{
			Type:    "info",
			Content: fmt.Sprintf("User %s not found", msg.Recipient),
		})
		return
	}
	err := recipientClient.Conn.WriteJSON(msg)
	if err != nil {
		log.Printf("error sending direct message to %s: %v", recipientClient.Username, err)
		closeConnection(recipientClient.Conn)
		delete(clients, recipientClient.Username)
	}
}
