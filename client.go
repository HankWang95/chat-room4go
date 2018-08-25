package chat_room4go

import (
	"bufio"
	"encoding/json"
	"github.com/smartwalle/xid"
	"log"
	"net"
	"time"
)

var clientMap = make(map[string]string)

type Client struct {
	identifier string
	conn       net.Conn
	hub        string
	postBox    chan *Message
}

func NewClient() *Client {
	var client = &Client{}
	identifier := string(xid.NewXIDWithTime(time.Now()))
	clientMap[identifier] = ""
	client.identifier = identifier
	return client
}

func (this *Client) Id() string {
	return this.identifier
}

func (this *Client) JoinHub(hub, serverAddr string) (postBox *chan *Message, err error) {
	if this.hub != "" {
		this.LeaveHub()
	}
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	this.conn = conn
	this.hub = hub
	this.sendMessage("", "", MSG_TYPE_SAY_HI)
	this.postBox = make(chan *Message, 10)
	go this.readMessage()
	clientMap[this.identifier] = this.hub
	return &this.postBox, nil
}

func (this *Client) SendBroadcast(msg string) {
	if this.hub == "" {
		log.Print("你没加入任何hub")
		return
	}
	this.sendMessage(msg, "", MSG_TYPE_BROADCAST)

}

func (this *Client) SendMessageTo(msg, to string) {
	if this.hub == "" {
		log.Print("你没加入任何hub")
		return
	}
	this.sendMessage(msg, to, MSG_TYPE_WHISPER)
}

func (this *Client) LeaveHub() {
	this.sendMessage("", "", MSG_TYPE_SAY_BYE)
	this.conn.Close()
	this.conn = nil
	close(this.postBox)
	this.hub = ""
	delete(clientMap, this.identifier)
}

func (this *Client) readMessage() {
	reader := bufio.NewReader(this.conn)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			log.Print(err)
			break
		}
		var msg = &Message{}
		err = json.Unmarshal(data, msg)
		if err != nil {
			log.Print(err)
			continue
		}
		this.postBox <- msg
	}
}

func (this *Client) sendMessage(msg, to string, mType int) {
	message := Message{
		From:    this.identifier,
		Hub:     this.hub,
		MsgType: mType,
		Msg:     msg,
		To:      to,
	}
	data, _ := json.Marshal(message)
	data = append(data, '\n')
	this.conn.Write(data)
}
