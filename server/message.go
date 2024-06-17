package main

import (
	"encoding/json"
	"log"
	"os"
)

type Message struct {
	Type      string `json:"type"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	Room      string `json:"room"`
	Recipient string `json:"recipient"`
	Password  string `json:"password"`
}

func loadRoomsFromJSON() {
	file, err := os.ReadFile(roomFile)
	if err != nil {
		if os.IsNotExist(err) {

			data := map[string][]map[string]interface{}{"rooms": {}}
			saveRoomsToJSON(data)
			return
		}
		log.Printf("Error reading %s: %v", roomFile, err)
		return
	}

	var data map[string][]map[string]interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Printf("Error unmarshaling %s: %v", roomFile, err)
		return
	}

	for _, room := range data["rooms"] {
		roomName := room["name"].(string)
		members := room["members"].([]interface{})

		rooms[roomName] = make(map[string]*Client)

		for _, member := range members {
			username := member.(string)
			if client, ok := clients[username]; ok {
				rooms[roomName][username] = client
				client.Rooms[roomName] = true
			}
		}
	}
}

func saveRoomsToJSON(data map[string][]map[string]interface{}) {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Printf("Error marshaling data to JSON: %v", err)
		return
	}

	err = os.WriteFile(roomFile, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing %s: %v", roomFile, err)
		return
	}
}

func loadUsersFromJSON() {
	file, err := os.ReadFile(usersFile)
	if err != nil {
		if os.IsNotExist(err) {

			data := make(map[string]string)
			saveUsersToJSON(data)
			return
		}
		log.Printf("Error reading %s: %v", usersFile, err)
		return
	}

	err = json.Unmarshal(file, &users)
	if err != nil {
		log.Printf("Error unmarshaling %s: %v", usersFile, err)
	}
}

func saveUsersToJSON(data map[string]string) {
	file, err := os.Create(usersFile)
	if err != nil {
		log.Printf("Error creating %s: %v", usersFile, err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		log.Printf("Error encoding %s: %v", usersFile, err)
	}
}

func generateRoomsJSON() map[string][]map[string]interface{} {
	roomsData := make([]map[string]interface{}, 0, len(rooms))
	for roomName, clients := range rooms {
		members := make([]string, 0, len(clients))
		for username := range clients {
			members = append(members, username)
		}
		roomData := map[string]interface{}{
			"name":    roomName,
			"members": members,
		}
		roomsData = append(roomsData, roomData)
	}
	return map[string][]map[string]interface{}{
		"rooms": roomsData,
	}
}
