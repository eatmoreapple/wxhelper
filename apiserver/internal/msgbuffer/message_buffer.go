package msgbuffer

import (
	"context"
	"errors"
	"github.com/eatmoreapple/env"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/go-redis/redis/v8"
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
	if add := env.Name("MSG_QUEUE_ADDR").String(); len(add) > 0 {
		// 先简单一点
		client := redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    add,
		})
		return NewRedisMessageBuffer(client, "")
	}
	return NewMemoryMessageBuffer(100)
}
