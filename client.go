package chat_room4go

import (
	"net"
)



type Client struct {
	identifier string
	conn *net.Conn
	hub  string
	send chan []byte
}

func NewClient(hub, serverAddr string) (client *Client, err error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	client.conn = &conn
	client.hub = hub
	client.send = make(chan []byte, 256)
	return client, nil
}
