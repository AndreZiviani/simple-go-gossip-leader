package main

import (
	"fmt"
	"io"

	"github.com/hashicorp/memberlist"
)

func Gossip(msgCh chan []byte, me string) *memberlist.Memberlist {

	d := &Delegate{
		msgCh: msgCh,
	}

	config := memberlist.DefaultLocalConfig()
	config.Name = me
	config.BindAddr = "0.0.0.0"
	config.BindPort = 3100
	config.AdvertisePort = config.BindPort
	config.Delegate = d
	config.LogOutput = io.Discard

	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	return list
}

func join(list *memberlist.Memberlist) {
	addr := "test-headless:3100"
	fmt.Printf("Joining cluster at %s\n", addr)
	_, err := list.Join([]string{addr})
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}
}

func broadcast(list *memberlist.Memberlist, msg Message) {
	for _, member := range list.Members() {
		if member.Name != list.LocalNode().Name {
			list.SendReliable(member, msg.Message())
		}
	}
}
