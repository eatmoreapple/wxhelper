package apiserver

import (
	"encoding/json"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"net/http"
)

// MessageHandler 消息处理者
type MessageHandler interface {
	ServeMessage(*Message)
}

// MessageHandlerFunc 以函数的方式实现MessageHandler
type MessageHandlerFunc func(*Message)

// ServeMessage 实现了MessageHandler
func (m MessageHandlerFunc) ServeMessage(msg *Message) { m(msg) }

// ReaderMessageHandler 从reader中解析 Message 并处理
type ReaderMessageHandler struct {
	Reader         io.Reader
	MessageHandler MessageHandler
}

func (r ReaderMessageHandler) Serve() error {
	var msg Message
	if err := json.NewDecoder(r.Reader).Decode(&msg); err != nil {
		log.Warn().Err(err).Msg("decode message failed")
		return err
	}
	log.Info().Msg("parse message successfully")
	r.MessageHandler.ServeMessage(&msg)
	return nil
}

type MessageListener interface {
	ListenAndServe(handler MessageHandler) error
}

type TCPMessageListener struct {
	Addr string
}

func (t *TCPMessageListener) ListenAndServe(messageHandler MessageHandler) error {
	listener, err := net.Listen("tcp", t.Addr)
	if err != nil {
		return err
	}
	defer func() { _ = listener.Close() }()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go t.processMessage(conn, messageHandler)
	}
}

func (t *TCPMessageListener) processMessage(conn net.Conn, messageHandler MessageHandler) {
	defer func() { _ = conn.Close() }()
	defer func() { _, _ = conn.Write([]byte("200 OK")) }()
	handler := ReaderMessageHandler{Reader: conn, MessageHandler: messageHandler}
	_ = handler.Serve()
}

type HTTPMessageListener struct {
	Addr    string
	handler MessageHandler
}

func (t *HTTPMessageListener) ListenAndServe(handler MessageHandler) error {
	t.handler = handler
	return http.ListenAndServe(t.Addr, t)
}

func (t *HTTPMessageListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() { _, _ = w.Write([]byte(`{"code": 0, "msg": "success"}`)) }()
	handler := ReaderMessageHandler{Reader: r.Body, MessageHandler: t.handler}
	_ = handler.Serve()
}
