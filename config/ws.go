package config

import (
	"encoding/json"
	"fmt"
	"go/chat/model"
	"go/chat/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
)

var SERVERID string = utils.EnvVariable("SERVERID")
var broadcast = make(chan *redis.Message)

type Message struct {
	Id        string
	Message   string `json:"msg"`
	Sender    string
	Receiver  string `json:"receiver,omitempty"`
	Group     bool   `json:"group"`
	GroupName string `json:"group_name,omitempty"`
}

var clients = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Wshandler(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	userID := c.Query("id")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	NewClient(userID, conn)
	model.SetUser(userID, SERVERID)
	ReceiveMessage(conn, userID)

}
func ReceiveMessage(conn *websocket.Conn, userID string) {
	for {
		_, msg, errCon := conn.ReadMessage()

		if errCon != nil {
			log.Println("Read Error:", errCon)
			conn.Close()
			break
		}
		var r Message
		if err := json.Unmarshal(msg, &r); err != nil {

			log.Println("Error: " + err.Error())
			return
		}
		r.Id = utils.GenerateKsuid()
		r.Sender = userID

		//find the server inwhich the receiver is connected
		serverId := model.GetServerId(r.Receiver)
		//send message to redis queue
		JsonData, err := json.Marshal(r)
		utils.CheckErr(err)
		model.SaveMessage(r.Id, r.Message, r.Sender, r.Receiver, r.Group, r.GroupName)
		Conn.Publish(Ctx, serverId, JsonData)
	}
}

func NewClient(ID string, conn *websocket.Conn) {
	clients[ID] = conn
	clients[ID].WriteMessage(websocket.TextMessage, []byte("ok"))
}

func Send() {
	for {
		msg := <-broadcast
		message := Message{}
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			panic(err)
		}
		JsonData, err := json.Marshal(message)
		utils.CheckErr(err)
		client := clients[message.Receiver]
		err = client.WriteMessage(websocket.TextMessage, []byte(JsonData))
		if err != nil {
			client.Close()
		}
	}
}
