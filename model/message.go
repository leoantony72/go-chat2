package model

import (
	"go/chat/database"
	"go/chat/utils"
)

func SaveMessagePrivateChat(id string, msg string, sender string, receiver string) {
	query := `INSERT INTO private_chat(id,msg,sender,receiver,timestamp)VALUES(?,?,?,?,toTimeStamp(now()))`
	err := database.ExecuteQuery(query, id, msg, sender, receiver)
	utils.CheckErr(err)
}

func SaveMessageGroupChat(id string, msg string, sender string, groupName string) {
	query := `INSERT INTO group_chat(id,msg,sender,group,timestamp)VALUES(?,?,?,?,toTimeStamp(now()))`
	err := database.ExecuteQuery(query, id, msg, sender,groupName)
	utils.CheckErr(err)
}
