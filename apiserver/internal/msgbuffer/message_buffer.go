package msgbuffer

import (
	"context"
	"errors"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	ErrTimeout = errors.New("get message timeout")
)

type MessageBuffer interface {
	Put(ctx context.Context, msg *Message) error
	Get(ctx context.Context, timeout time.Duration) (*Message, error)
}

type MemoryMessageBuffer struct {
	msgCH chan *Message
}

func (m *MemoryMessageBuffer) Put(ctx context.Context, msg *Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m.msgCH <- msg:
		log.Info().Interface("message", msg).Msg("put message to buffer")
	default:
		log.Warn().Msg("message buffer is full")
	}
	return nil
}

func (m *MemoryMessageBuffer) Get(ctx context.Context, timeout time.Duration) (*Message, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		log.Info().Msg("get message from buffer")
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
			return nil, ErrTimeout
		case msg := <-m.msgCH:
			return msg, nil
		}
	}
}

func NewMemoryMessageBuffer(size int) MessageBuffer {
	return &MemoryMessageBuffer{msgCH: make(chan *Message, size)}
}

func Default() MessageBuffer {
	return NewMemoryMessageBuffer(100)
}
