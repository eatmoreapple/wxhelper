package taskpool

import (
	"context"
	"golang.org/x/sync/semaphore"
)

var _semaphore = semaphore.NewWeighted(100)

func Do(ctx context.Context, handler func()) error {
	if err := _semaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	go func() {
		defer _semaphore.Release(1)
		handler()
	}()
	return nil
}
