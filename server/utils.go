package main

import (
	"fmt"
	"strings"
)

func notifyRoomUpdate(roomName string) {
	if clientsInRoom, ok := rooms[roomName]; ok {
		members := make([]string, 0, len(clientsInRoom))
		for username := range clientsInRoom {
			members = append(members, username)
		}
		updateMessage := Message{
			Type:    "room_update",
			Room:    roomName,
			Content: fmt.Sprintf("Room %s members: %s", roomName, strings.Join(members, ", ")),
		}
		for _, client := range clientsInRoom {
			client.Conn.WriteJSON(updateMessage)
		}
	}
}
