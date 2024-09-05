package main

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/memberlist"
)

type Delegate struct {
	msgCh chan []byte
}

func (d *Delegate) NotifyMsg(msg []byte) {
	d.msgCh <- msg
}
func (d *Delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}
func (d *Delegate) NodeMeta(limit int) []byte {
	// not use, noop
	return []byte("")
}
func (d *Delegate) LocalState(join bool) []byte {
	// not use, noop
	return []byte("")
}
func (d *Delegate) MergeRemoteState(buf []byte, join bool) {
	// not use
}

type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m Message) Invalidates(other memberlist.Broadcast) bool {
	return false
}
func (m Message) Finished() {
	// nop
}
func (m Message) Message() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("")
	}
	return data
}
func (m Message) String() string {
	return fmt.Sprintf("Message{Key: %s, Value: %s}", m.Key, m.Value)
}

func ParseMessage(data []byte) (Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	return m, err
}
