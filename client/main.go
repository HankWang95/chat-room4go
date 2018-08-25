package main

import (
	"github.com/HankWang95/chat-room4go"
	"time"
)

func main() {
	client1 := chat_room4go.NewClient()

	client2 := chat_room4go.NewClient()
	client1.JoinHub("haha", ":7777")
	client2.JoinHub("haha", ":7777")
	client1.SendBroadcast("broadcast")
	time.Sleep(time.Second*3)
	client1.SendBroadcast("broadcast")
	client1.SendBroadcast("broadcast")
	time.Sleep(time.Second*3)

	client1.SendBroadcast("broadcast")

	client1.SendMessageTo("hi" ,client2.Id())
}
