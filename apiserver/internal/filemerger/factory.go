package filemerger

import (
	"errors"
	"github.com/eatmoreapple/wxhelper/apiserver/internal/filemerger/internal"
	"strings"
	_ "unsafe"
)

// Factory is an interface that defines a method for creating a new FileMerger.
type Factory interface {
	// New creates a new FileMerger.
	// It takes a key string as a parameter.
	// It returns a new FileMerger and an error if the operation fails.
	New(key string) (FileMerger, error)
}

// do not import "github.com/eatmoreapple/wxhelper/internal/wxclient" directly
//
//go:linkname tempDir github.com/eatmoreapple/wxhelper/internal/wxclient.tempDir
func tempDir() string

// localFileMergerFactory is a struct that holds a Redis client.
type localFileMergerFactory struct {
	cache internal.Cache
}

// New is a method that creates a new instance of localFileMerger.
func (l *localFileMergerFactory) New(key string) (FileMerger, error) {
	items := strings.Split(key, ":")
	if len(items) < 2 {
		return nil, errors.New("invalid key")
	}
	length := len(items)
	filename, fileHash := strings.Join(items[:length-1], ":"), items[length-1]
	return &localFileMerger{
		cache:    l.cache,
		filename: filename,
		fileHash: fileHash,
		tempDir:  tempDir(),
	}, nil
}

// DefaultFactory is a function that creates a new Factory.
// It reads the REDIS_ADDR environment variable to configure the Redis client.
// It panics if the REDIS_ADDR environment variable is not set.
// It returns a Factory.
func DefaultFactory() Factory {
	cache := internal.CacheFromEnv()
	return &localFileMergerFactory{cache: cache}
}
