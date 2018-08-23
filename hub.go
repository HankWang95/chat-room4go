package chat_room4go

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)


const (
	MSG_TYPE_BROADCAST = 1000 * (iota + 1)
	MSG_TYPE_WHISPER
)

type Hub struct {
	identifier     string
	connectionsMap map[string]net.Conn
	message        chan *Message
}

func NewHub(identifier string) *Hub {
	var hub *Hub
	hub.identifier = identifier
	hub.message = make(chan *Message)
	hub.connectionsMap = make(map[string]net.Conn)
	hubList[hub.identifier] = hub
	return hub
}

func (this *Hub) broadcast(msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	for _, connection := range this.connectionsMap {
		connection.Write(data)
	}
}

func (this *Hub) whisper(msg *Message) {
	receiver := this.connectionsMap[msg.To]
	if receiver == nil {
		log.Fatal("该用户不在线")
	}
	data, _ := json.Marshal(msg)

	receiver.Write(data)
}

func (this *Hub) addConnection(conn net.Conn, identifier string) {
	if identifier == "" {
		log.Fatal("identifier is not exist")
	}
	this.connectionsMap[identifier] = conn

	go this.handleConnect(conn, identifier)
}

func (this *Hub) handleMessage() {
	for {
		msg := <-this.message
		switch msg.MsgType {
		case MSG_TYPE_BROADCAST:
			this.broadcast(msg)
		case MSG_TYPE_WHISPER:
			this.whisper(msg)
		}
	}
}

func (this *Hub) handleConnect(conn net.Conn, identifier string) {
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		var msg *Message
		json.Unmarshal(data, msg)
		this.message <- msg
	}

	// When we encouter `err` reading, send this
	// connection to `deadConnections` for removal.
	//
	conn.Close()
}
