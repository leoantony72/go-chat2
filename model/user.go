package model

import "go/chat/database"

func CreateUser(userId string, userName string) {
	query := `INSERT INTO users(id,username)VALUES(?,?)`

	database.ExecuteQuery(query, userId, userName)
}

func CheckUserExist(userId string) (string, string) {
	query := `SELECT id,username FROM users WHERE username = ?`

	//check the database for id
	ID, Username := database.CheckUserExist(query, userId)
	return ID, Username
}
