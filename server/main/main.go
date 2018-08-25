package main

import "github.com/HankWang95/chat-room4go"

func main()  {
	chat_room4go.NewServer(":7777")
	chat_room4go.NewHub("haha")
	select {

	}
}