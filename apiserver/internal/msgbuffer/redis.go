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
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.client.LPush(ctx, r.queue, data).Err()
}

func (r RedisMessageBuffer) Get(ctx context.Context, timeout time.Duration) (*Message, error) {
	msgs, err := r.client.BRPop(ctx, timeout, r.queue).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNoMessage
	}
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, ErrNoMessage
	}
	if len(msgs) != 2 {
		// unreachable, but just in case
		return nil, errors.New("invalid message")
	}
	var msg Message
	if err = json.Unmarshal([]byte(msgs[1]), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func NewRedisMessageBuffer(client *redis.Client, queue string) *RedisMessageBuffer {
	if queue == "" {
		queue = "wechat:message:queue"
	}
	return &RedisMessageBuffer{client: client, queue: queue}
}
