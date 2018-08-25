package chat_room4go

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)



type Message struct {
	From    string `json:"from"`
	Hub     string `json:"hub"`
	MsgType int    `json:"msg_type"`
	Msg     string `json:"msg"`
	To      string `json:"to"`
}

func NewServer(serverAddr string)  {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal("tcp 聊天室连接失败")
	}
	go chatRoomListener(listener)
}

func chatRoomListener(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("chat-room hub tcp-conn error!")
		}
		reader := bufio.NewReader(conn)
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Print(err)
			continue
		}
		msg := &Message{}
		json.Unmarshal(data, msg)

		if hub ,ok := hubList[msg.Hub]; ok {
			hub.addConnection(conn, msg.From)
		} else{
			msg.From = "系统通知"
			msg.Msg = "该hub不存在"
			data, _ := json.Marshal(msg)
			conn.Write(data)
			conn.Close()
		}
	}
}
