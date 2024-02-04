package msgbuffer

import (
	"context"
	"errors"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"time"
)

var (
	// ErrNoMessage is returned when no message is found in the buffer.
	ErrNoMessage = errors.New("no message found")
)

type MessageBuffer interface {
	// Put adds a message to the buffer.
	Put(ctx context.Context, msg *Message) error

	// Get retrieves a message from the buffer.
	// If no message is available, it will return ErrNoMessage.
	Get(ctx context.Context, timeout time.Duration) (*Message, error)
}

func Default() MessageBuffer {
	return NewMemoryMessageBuffer(100)
}
