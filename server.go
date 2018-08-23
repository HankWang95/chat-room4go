package chat_room4go

import (
	"net"
	"log"
	"bufio"
	"encoding/json"
)


type chatServer struct {
	net.Listener
}

var hubList = make(map[string]*Hub)

type Message struct {
	From    string `json:"from"`
	Hub     string `json:"hub"`
	MsgType int    `json:"msg_type"`
	Msg     string `json:"msg"`
	To      string `json:"to"`
}

func NewServer(serverAddr string) *chatServer {
	server, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal("tcp 聊天室连接失败")
	}
	var chatServer = &chatServer{server}
	return chatServer
}

func (this *chatServer) Listen() {
	for {
		conn, err := this.Accept()
		if err != nil {
			log.Fatal("chat-room hub tcp-conn error!")
		}
		reader := bufio.NewScanner(conn)
		data := reader.Bytes()
		msg := &Message{}
		json.Unmarshal(data, msg)

		var hub = &Hub{}

		for hubId := range hubList {
			if hubId == msg.Hub {
				hub = hubList[hubId]
				hub.addConnection(conn, msg.From)
				break
			}
		}
	}
}

