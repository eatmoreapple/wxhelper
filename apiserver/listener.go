package apiserver

import (
	"bytes"
	"context"
	"github.com/eatmoreapple/wxhelper/apiserver/internal/msgbuffer"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"io"
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
		go func(coon net.Conn) {
			if err = m.serve(conn); err != nil {
				log.Error().Err(err).Msg("failed to serve")
			}
		}(conn)
	}
}

func (m *MessageListener) serve(coon net.Conn) error {
	log.Info().Msg("receive new message")
	defer func() { _ = coon.Close() }()
	defer func() { _, _ = coon.Write([]byte("200 OK")) }()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, coon); err != nil {
		return err
	}
	var msg Message
	if err := json.NewDecoder(&buf).Decode(&msg); err != nil {
		return err
	}
	log.Info().Msg("parse message successfully")
	return m.MessageBuffer.Put(context.TODO(), &msg)
}

func NewMessageListener(msgBuffer msgbuffer.MessageBuffer) *MessageListener {
	return &MessageListener{
		MessageBuffer: msgBuffer,
	}
}
