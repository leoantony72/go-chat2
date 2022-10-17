package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

var Ctx = context.Background()
var Conn redis.Client

func NPool() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	Conn = *rdb
}

func PubSub() {
	subscriber := Conn.Subscribe(Ctx, "server1")
	for {
		msg, err := subscriber.ReceiveMessage(Ctx)
		if err != nil {
			panic(err)
		}

		fmt.Printf("message from pub/sub : %v", msg.Payload)
		broadcast <- msg

		// ...
	}
}
