package model

import (
	"go/chat/database"
)

func CreateRoom(id string, roomName string) {
	query := `INSERT INTO room (id,room_name) VALUES (?,?)`
	database.ExecuteQuery(query, id, roomName)
}

func JoinRoom(roomName string, username string) {
	query := `INSERT INTO room_members(room_name,username)VALUES(?,?)`

	database.ExecuteQuery(query, roomName, username)
}
func GetMembers(groupName string) []string {
	var username string
	members := []string{}
	query := "SELECT username FROM room_members WHERE room_name = ?"
	data := database.Connection.Session.Query(query, groupName).Iter()
	for data.Scan(&username) {
		members = append(members, username)
	}

	return members
}
