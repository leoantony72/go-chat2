package model

import "go/chat/database"

func SaveMessage(id string, msg string, sender string, receiver string, group bool, groupName string) {
	query := `INSERT INTO message(id,msg,sender,receiver,isgroup,group)VALUES(?,?,?,?,?,?)`
	database.ExecuteQuery(query, id, msg, sender, receiver, group, groupName)
}
