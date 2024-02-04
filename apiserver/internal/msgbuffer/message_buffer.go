package msgbuffer

import (
	"context"
	"errors"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"time"
)

var (
	ErrTimeout = errors.New("get message timeout")
)

type MessageBuffer interface {
	Put(ctx context.Context, msg *Message) error
	Get(ctx context.Context, timeout time.Duration) (*Message, error)
}

func Default() MessageBuffer {
	return NewMemoryMessageBuffer(100)
}
