// Package filemerger provides interfaces and functions for merging files.
package filemerger

import (
	"context"
	"github.com/eatmoreapple/env"
	"github.com/go-redis/redis/v8"
)

// FileMerger is an interface that defines methods for adding and merging files.
type FileMerger interface {
	// Add adds a file to the merger.
	// It takes a context and a file string as parameters.
	// It returns an error if the operation fails.
	Add(ctx context.Context, file string) error

	// Merge merges all added files.
	// It takes a context as a parameter.
	// It returns a string representing the merged file and an error if the operation fails.
	Merge(ctx context.Context) (string, error)
}

// Factory is an interface that defines a method for creating a new FileMerger.
type Factory interface {
	// New creates a new FileMerger.
	// It takes a key string as a parameter.
	// It returns a new FileMerger and an error if the operation fails.
	New(key string) (FileMerger, error)
}

// DefaultFactory is a function that creates a new Factory.
// It reads the REDIS_ADDR environment variable to configure the Redis client.
// It panics if the REDIS_ADDR environment variable is not set.
// It returns a Factory.
func DefaultFactory() Factory {
	addr := env.Name("REDIS_ADDR").String()
	if len(addr) == 0 {
		panic("REDIS_ADDR is required")
	}
	return &localFileMergerFactory{client: redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     addr,
		PoolSize: 10,
	})}
}
