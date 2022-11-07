package config

import (
	"encoding/json"
	"fmt"
	"go/chat/model"
	"go/chat/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
)

var SERVERID string = utils.EnvVariable("SERVERID")
var broadcast = make(chan *redis.Message)

type Message struct {
	Id           string
	Message      string   `json:"msg"`
	Sender       string   `json:"sender"`
	Receiver     string   `json:"receiver,omitempty"`
	Group        bool     `json:"is_group"`
	GroupName    string   `json:"group_name,omitempty"`
	GroupMembers []string `json:",omitempty"`
	ServerId     string   `json:",omitempty"`
}
type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
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
	_, username := model.CheckUserExist(userID)
	if username == "" {
		Closews("Authentication failed - invalid username", conn)
		return
	}

	NewClient(userID, conn)
	ReceiveMessage(conn, userID)

}
func ReceiveMessage(conn *websocket.Conn, userID string) {
	for {
		_, msg, errCon := conn.ReadMessage()

		if errCon != nil {
			log.Println("Read Error:", errCon)
			break
		}
		var r Message
		if err := json.Unmarshal(msg, &r); err != nil {
			log.Println("Error: " + err.Error())
			MsgFailed(conn)
			continue
		}
		r.Id = utils.GenerateKsuid()
		r.Sender = userID
		err := r.Validate()
		if err != nil {
			b, _ := json.Marshal(err)
			conn.WriteMessage(websocket.TextMessage, b)
			continue
		}
		if r.Group {
			/*
				@send the message to the db
				@fetch the members of the group
				@fetch the members connected server
				@send message to the topic of the servers
			*/
			model.SaveMessageGroupChat(r.Id, r.Message, r.Sender, r.GroupName)
			members := model.GetMembers(r.GroupName)
			// servers := []string{}
			servers := make(map[string][]string)
			for _, member := range members {
				serverId := model.GetServerId(member)
				servers[serverId] = append(servers[serverId], member)
			}
			// fmt.Println(servers)
			for key, element := range servers {
				r.ServerId = key
				r.GroupMembers = element
				JsonData, err := json.Marshal(r)
				utils.CheckErr(err)
				fmt.Println("redis key", key)
				Conn.Publish(Ctx, key, JsonData)

			}
			continue

		}
		model.SaveMessagePrivateChat(r.Id, r.Message, r.Sender, r.Receiver)
		//find the server inwhich the receiver is connected
		serverId := model.GetServerId(r.Receiver)
		//send message to redis queue
		JsonData, err := json.Marshal(r)

		utils.CheckErr(err)
		Conn.Publish(Ctx, serverId, JsonData)
	}
	cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Connection Closing")
	if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
		utils.CheckErr(err)
	}
	conn.Close()

}

func NewClient(userId string, conn *websocket.Conn) {
	model.SetUser(userId, SERVERID)
	clients[userId] = conn
	clients[userId].WriteMessage(websocket.TextMessage, []byte("ok"))
}

func Send() {
	for {
		msg := <-broadcast
		message := Message{}
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			panic(err)
		}
		if message.Group {
			groupMessage(message)
			continue
		}
		client := clients[message.Receiver]
		if client == nil {
			fmt.Println("Receiver offline")
			continue
		}
		privateMessage(message, client)
	}
}

func groupMessage(message Message) {
	jsonRes := Message{}
	for _, member := range message.GroupMembers {
		client := clients[member]
		if client == nil {
			fmt.Println("Receiver offline group")
			continue
		}
		jsonRes.Id = message.Id
		jsonRes.Sender = message.Sender
		jsonRes.Message = message.Message
		jsonRes.Group = message.Group
		jsonRes.GroupName = message.GroupName
		JsonData, err := json.Marshal(jsonRes)
		utils.CheckErr(err)
		err = client.WriteMessage(websocket.TextMessage, []byte(JsonData))
		if err != nil {
			delete(clients, message.Receiver)
			client.Close()
		}
	}
}
func privateMessage(message Message, client *websocket.Conn) {
	JsonData, err := json.Marshal(message)
	utils.CheckErr(err)
	err = client.WriteMessage(websocket.TextMessage, []byte(JsonData))
	if err != nil {
		delete(clients, message.Receiver)
		client.Close()
	}
}

func Closews(msg string, conn *websocket.Conn) {
	cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, msg)
	if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
		utils.CheckErr(err)
	}
	conn.Close()
}
func MsgFailed(conn *websocket.Conn) {
	msg := `{"message":"Failed to send message"}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		utils.CheckErr(err)
	}
}

func (m Message) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Message,
			validation.Required.Error("msg Field is required"),
			validation.NotNil.Error("msg field cannot be empty"),
			validation.Length(1, 1000).Error("character length should be between 1 and 1000"),
			// validation.Empty.Error("msg field is required"),
			// validation.NotNil.Error("msg field is required"),
		),
		validation.Field(&m.Group,
			// validation.Required.Error("is_group Field is required"),
			validation.NotNil.Error("is_group field cannot be empty"),
		// validation.Empty.Error("msg field is required"),
		// validation.NotNil.Error("msg field is required"),
		),
		validation.Field(&m.GroupName,
			validation.Length(1, 25).Error("character length should be between 1 and 25"),
			validation.When(m.Group, validation.Required.Error("Group_name is required")),
		),
		validation.Field(&m.Receiver,
			validation.When(m.Group, validation.Empty).Else(validation.Required.Error("reciever Field is required"), validation.NotNil.Error("receiver field cannot be empty")),
		),
	)

}
