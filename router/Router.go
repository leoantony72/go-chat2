package router

import (
	"go/chat/config"
	"go/chat/controllers"
	"go/chat/database"
	"go/chat/utils"

	"github.com/gin-gonic/gin"
)

var PORT string = utils.EnvVariable("PORT")

func StartServer() {
	router := gin.Default()
	// config.NewWebSocketServer()
	database.SetupDBconnection()
	config.NPool()
	go config.PubSub()
	go config.Send()

	router.GET("/", home)
	router.GET("/chat", test)                      //websocket_connection
	router.POST("/room", controllers.Createroom)   //@ create room
	router.POST("/joinroom", controllers.JoinRoom) //@ join room *http
	router.POST("/user", controllers.CreateUser)   //@ Create User
	router.POST("/login", controllers.LoginUser)   //@Login User
	router.Run(":" + PORT)
}
func test(c *gin.Context) {
	// ctx := context.Background()
	config.Wshandler(c.Writer, c.Request, c)

}

func home(c *gin.Context) {

	c.JSON(200, gin.H{"message": "server started"})
}
