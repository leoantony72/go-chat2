package controllers

import (
	"fmt"
	"go/chat/database"
	"go/chat/utils"

	"github.com/gin-gonic/gin"
	// "github.com/gorilla/websocket"
)

type User struct {
	// conn     *websocket.Conn
	Id       string `json:"user_id"`
	Username string `json:"username"`
}
type LoginRequest struct {
	Id string `json:"id"`
}

func CreateUser(c *gin.Context) {
	user := &User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println(err)
	}

	Id := utils.GenerateKsuid()
	query := `INSERT INTO users (id,username) VALUES(?,?)`

	database.ExecuteQuery(query, Id, user.Username)

	c.JSON(201, gin.H{"message": "User Created"})
}

func LoginUser(c *gin.Context) {
	//get id from form
	user := &LoginRequest{}
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println(err)
	}
	query := `SELECT id,username FROM users WHERE id = ?`

	//check the database for id
	ID, Username := database.CheckUserExist(query, user.Id)

	//if exist return a cookie containing the ID
	if ID == "" {
		c.JSON(200, gin.H{"error": "invalid user"})
		return
	}
	c.SetCookie("uid", ID, 36000, "/", "localhost", false, true)
	c.JSON(200, gin.H{"id": ID, "name": Username})
}
