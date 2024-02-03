package apiserver

import (
	"github.com/eatmoreapple/wxhelper/apiserver/internal/msgbuffer"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type MessageListener interface {
	http.Handler
}

type messageListener struct {
	msgbuffer.MessageBuffer
}

func (m *messageListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("receive message")
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Error().Err(err).Msg("decode message failed")
		return
	}
	log.Info().Interface("message", msg).Msg("receive message")
	if err := m.MessageBuffer.Put(r.Context(), &msg); err != nil {
		log.Error().Err(err).Msg("put message to buffer failed")
	}
	// {"code": 0, "msg": "success"}
	_, _ = w.Write([]byte(`{"code": 0, "msg": "success"}`))
}

func NewMessageListener(msgBuffer msgbuffer.MessageBuffer) MessageListener {
	return &messageListener{
		MessageBuffer: msgBuffer,
	}
}
