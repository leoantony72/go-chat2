package model

import (
	"go/chat/database"
	"go/chat/utils"
)

func SetUser(userID string, SERVERID string) {
	query := "INSERT INTO user_mapping(username,server_id)VALUES(?,?)"
	err:=database.ExecuteQuery(query, userID, SERVERID)
	utils.CheckErr(err)
}

func GetServerId(userID string) string {
	var ServerId string
	// var ids []string
	query := `SELECT server_id FROM user_mapping WHERE username =?`
	iter := database.Connection.Session.Query(query, userID).Iter()
	iter.Scan(&ServerId)

	return ServerId
}
