package main

import (
	"fmt"
	"github.com/benbpyle/readings-trasmitter/models"
	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	psc.Subscribe("c1")

	for {
		switch n := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("Message: %s %s\n", n.Channel, n.Data)
			r := models.NewTempReading(n.Data)
			if r != nil {
				models.Insert(r)
			}
		case redis.Subscription:
			fmt.Printf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
			if n.Count == 0 {
				return
			}
		case error:
			fmt.Printf("error: %v\n", n)
			return
		}
	}
}
