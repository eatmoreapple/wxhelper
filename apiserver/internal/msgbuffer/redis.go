package msgbuffer

import (
	"context"
	"encoding/json"
	"errors"
	. "github.com/eatmoreapple/wxhelper/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisMessageBuffer struct {
	client *redis.Client
	queue  string
}

func (r RedisMessageBuffer) Put(ctx context.Context, msg *Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return r.client.LPush(ctx, r.queue, msg).Err()
}

func (r RedisMessageBuffer) Get(ctx context.Context, timeout time.Duration) (*Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	msgs, err := r.client.BRPop(ctx, timeout, r.queue).Result()
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, ErrTimeout
	}
	if len(msgs) != 2 {
		return nil, errors.New("invalid message")
	}
	var msg Message
	if err = json.Unmarshal([]byte(msgs[1]), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func NewRedisMessageBuffer(client *redis.Client) *RedisMessageBuffer {
	return &RedisMessageBuffer{client: client}
}
