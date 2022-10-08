package config

import (
	// "fmt"

	"fmt"
	"log"
	"time"

	// "go/chat/utils"
	// "log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// type User struct {
// 	conn *websocket.Conn
// 	Id   string
// 	Send chan []byte
// }

var clients = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// func NewWebSocketServer() *ws {
// 	return &ws{
// 		users: make(map[string]User),
// 	}
// }

func Wshandler(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	ID := c.Query("id")
	// utils.CheckErr(err)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	NewClient(ID, conn)

	fmt.Println(clients)

}

func NewClient(ID string, conn *websocket.Conn) {
	// user := &User{
	// 	Id:   ID,
	// 	conn: conn,
	// 	Send: make(chan []byte),
	// }
	clients[ID] = conn
	// fmt.Println(ws.users[ID])
	// fmt.Println(ws.users[ID].Id)
	clients[ID].WriteMessage(websocket.TextMessage, []byte("hello"))

	// fmt.Println(ws)
}

func Echo() {
	for {
		time.Sleep(time.Second)
		// ws.users["2FhfPK3IvyicuLq9MxfuGFEK2eo"].conn.WriteMessage(websocket.TextMessage, []byte("hello"))
		// // send to every client that is currently connected
		for key, client := range clients {
			fmt.Println(key)
			err := client.WriteMessage(websocket.TextMessage, []byte("hello"))
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, key)
			}
		}
	}
}
