package wxhelper

import (
	"fmt"
	"testing"
)

func TestClient_CheckLogin(t *testing.T) {
	bot := New()
	bot.MessageHandler = func(msg *Message) {
		fmt.Println(msg.Content)
		if msg.Content == "ping" {
			fmt.Println(msg.ReplyText("pong"))
		}
	}
	fmt.Println(bot.Run())
}
