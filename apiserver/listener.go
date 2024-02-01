package apiserver

import (
	"github.com/eatmoreapple/wxhelper/apiserver/internal/msgbuffer"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/goccy/go-json"
	"net/http"
)

type MessageListener interface {
	http.Handler
}

type messageListener struct {
	msgbuffer.MessageBuffer
}

func (m messageListener) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		return
	}
	_ = m.MessageBuffer.Put(r.Context(), &msg)
}

func NewMessageListener(msgBuffer msgbuffer.MessageBuffer) MessageListener {
	return messageListener{
		MessageBuffer: msgBuffer,
	}
}
