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
	Id        string
	Message   string `json:"msg"`
	Sender    string
	Receiver  string `json:"receiver,omitempty"`
	Group     bool   `json:"is_group"`
	GroupName string `json:"group_name,omitempty"`
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
		fmt.Printf("Username not found")
		Closews("Authentication failed - invalid username", conn)
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
			// fmt.Println(string(b))
			conn.WriteMessage(websocket.TextMessage, b)
			continue
		}
		if r.Group {
			/*
				@send the message to the db
				@fetch the members of the group
				@fetch the members connected servers
				@send message topic of the servers
			*/

			model.SaveMessage(r.Message, r.Sender, r.Receiver, r.Group, r.GroupName)
			members := model.GetMembers(r.GroupName)
			servers := []string{}
			for _, mem := range members {
				serverId := model.GetServerId(mem)
				fmt.Println(serverId)
				added := isAdded(servers, serverId)
				if !added {
					servers = append(servers, serverId)
				}
			}
			fmt.Println(servers)

		}
		//find the server inwhich the receiver is connected
		serverId := model.GetServerId(r.Receiver)
		//send message to redis queue
		JsonData, err := json.Marshal(r)
		utils.CheckErr(err)
		model.SaveMessage(r.Message, r.Sender, r.Receiver, r.Group, r.GroupName)
		Conn.Publish(Ctx, serverId, JsonData)
	}
	cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Connection Closing")
	if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
		utils.CheckErr(err)
	}
	conn.Close()

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
		if client == nil {
			/*
				Notification service handles the
				message when the reciever is offline
				or something unexpected happens
			*/
			fmt.Println("Receiver offline")
			continue
		}
		fmt.Println("here")
		err = client.WriteMessage(websocket.TextMessage, []byte(JsonData))
		if err != nil {
			delete(clients, message.Receiver)
			client.Close()
		}
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
		validation.Field(&m.Group),
		// validation.Required.Error("is_group Field is required"),
		// validation.NotNil.Error("is_group field cannot be empty"),
		// validation.Empty.Error("msg field is required"),
		// validation.NotNil.Error("msg field is required"),
		// ),
		validation.Field(&m.GroupName,
			validation.Length(1, 25).Error("character length should be between 1 and 25"),
			validation.When(m.Group, validation.Required.Error("Group_name is required")),
		),
		validation.Field(&m.Receiver,
			validation.When(!m.Group, validation.Required.Error("reciever Field is required"), validation.NotNil.Error("receiver field cannot be empty")).Else(validation.Empty),
		),
	)

}

func Closews(msg string, conn *websocket.Conn) {
	cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, msg)
	if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
		utils.CheckErr(err)
	}
	conn.Close()
}
func MsgFailed(conn *websocket.Conn) {
	msg := "Failed to send message"
	if err := conn.WriteMessage(websocket.CloseMessage, []byte(msg)); err != nil {
		utils.CheckErr(err)
	}
}
func isAdded(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
