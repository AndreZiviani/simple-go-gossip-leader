package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	msgCh := make(chan []byte)
	me := os.Getenv("HOSTNAME")

	list := Gossip(msgCh, me)

	join(list)

	le := leader(context.Background(), me)
	broadcaster := time.NewTicker(5 * time.Second)
	for {
		select {
		case data := <-msgCh:
			msg, err := ParseMessage(data)
			if err != nil {
				fmt.Println("Failed to parse message: " + err.Error())
				continue
			}
			fmt.Printf("Received message: %s\n", msg)

		case <-broadcaster.C:
			if le.IsLeader() {
				fmt.Println("Sending msg to all nodes")
				msg := Message{Key: "from node " + me, Value: "Hello Mr. Anderson, I'm the leader!"}
				broadcast(list, msg)
			}
		}
	}
}
