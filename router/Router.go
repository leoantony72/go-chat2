package router

import (
	"go/chat/config"
	"go/chat/controllers"
	"go/chat/database"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	// config.NewWebSocketServer()
	database.SetupDBconnection()

	router.GET("/")
	router.GET("/chat", test)                      //websocket_connection
	router.POST("/room", controllers.Createroom)   //@ create room
	router.POST("/joinroom", controllers.JoinRoom) //@ join room *http
	router.POST("/user", controllers.CreateUser)   //@ Create User
	router.POST("/login", controllers.LoginUser)   //@Login User
	router.Run("localhost:6300")
}
func test(c *gin.Context) {

	go config.Echo()
	// ctx := context.Background()
	config.Wshandler(c.Writer, c.Request, c)

}
