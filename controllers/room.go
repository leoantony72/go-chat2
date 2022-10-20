//create room after authentication
package controllers

import (
	// "fmt"

	"fmt"
	"go/chat/model"
	"go/chat/utils"

	"github.com/gin-gonic/gin"
)

type Room struct {
	Id   string `json:"room_id"`
	Name string `json:"name"`
	User string `json:"user_id"`
}

func Createroom(c *gin.Context) {
	newRoom := Room{}
	if err := c.ShouldBindJSON(&newRoom); err != nil {
		fmt.Println(err)
	}

	id := utils.GenerateKsuid()
	model.CreateRoom(id, newRoom.Name)
	c.JSON(200, gin.H{"mess": "done"})
}

func JoinRoom(c *gin.Context) {

	newRoom := Room{}
	if err := c.ShouldBindJSON(&newRoom); err != nil {
		fmt.Println(err)
	}

	model.JoinRoom(newRoom.Id, newRoom.User)
	c.JSON(200, gin.H{"mess": "Room Joined"})

}
