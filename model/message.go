package model

import "go/chat/database"

func SaveMessage(msg string, sender string, receiver string, group bool, groupName string) {
	query := `INSERT INTO message(msg,sender,receiver,isgroup,group)VALUES(?,?,?,?,?)`
	database.ExecuteQuery(query, msg, sender, receiver, group, groupName)
}
