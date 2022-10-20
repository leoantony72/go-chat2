package model

import "go/chat/database"

func CreateRoom(id string, roomName string) {
	query := `INSERT INTO room (id,room_name) VALUES (?,?)`
	database.ExecuteQuery(query, id, roomName)
}

func JoinRoom(roomId string, userId string) {
	query := `INSERT INTO room_members (room_id,user_id) VALUES (?,?)`

	database.ExecuteQuery(query, roomId, userId)
}
