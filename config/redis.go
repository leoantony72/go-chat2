package config

import (
	"context"
	"fmt"

	// "fmt"
	"go/chat/utils"

	"github.com/go-redis/redis/v9"
)

var Ctx = context.Background()
var Conn redis.Client

func NPool() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	Conn = *rdb
}

func PubSub() {
	SERVERID := utils.EnvVariable("SERVERID")
	fmt.Println(SERVERID)
	subscriber := Conn.Subscribe(Ctx, SERVERID)
	for {
		msg, err := subscriber.ReceiveMessage(Ctx)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("message from pub/sub : %v", msg.Payload)
		broadcast <- msg

		// ...
	}
}
