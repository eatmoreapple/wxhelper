package apiserver

import (
	"context"
	"github.com/eatmoreapple/wxhelper/apiserver/internal/msgbuffer"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"net"
)

type MessageListener struct {
	msgbuffer.MessageBuffer
}

func (m *MessageListener) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Info().Msg("listening on " + addr)
	defer func() { _ = listener.Close() }()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go m.serve(conn)
	}
}

func (m *MessageListener) serve(coon net.Conn) {
	log.Info().Msg("receive new message")
	defer func() { _ = coon.Close() }()
	defer func() { _, _ = coon.Write([]byte("200 OK")) }()
	var msg Message
	if err := json.NewDecoder(coon).Decode(&msg); err != nil {
		log.Warn().Err(err).Msg("decode message failed")
		return
	}
	log.Info().Msg("parse message successfully")
	if err := m.MessageBuffer.Put(context.TODO(), &msg); err != nil {
		log.Error().Err(err).Msg("put message failed")
	}
}

func NewMessageListener(msgBuffer msgbuffer.MessageBuffer) *MessageListener {
	return &MessageListener{
		MessageBuffer: msgBuffer,
	}
}
