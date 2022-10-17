package config

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

func NPool() redis.Conn {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Panic(err)
	}

	return conn
}
