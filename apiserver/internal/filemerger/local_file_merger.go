package filemerger

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/go-redis/redis/v8"
	"io"
	"os"
	"path/filepath"
	"strings"
	_ "unsafe"
)

// redisLocalFileMerger is a struct that holds a Redis client, filename and fileHash.
type redisLocalFileMerger struct {
	client   *redis.Client
	filename string
	fileHash string
	tempDir  string
}

// Add is a method that adds a file to a Redis list.
func (r *redisLocalFileMerger) Add(ctx context.Context, file string) error {
	return r.client.RPush(ctx, r.fileHash, file).Err()
}

// Merge is a method that merges all files in a Redis list and checks if the hash of the merged file matches the fileHash.
func (r *redisLocalFileMerger) Merge(ctx context.Context) (string, error) {
	// try to get all files from redis
	files, err := r.client.LRange(ctx, r.fileHash, 0, -1).Result()
	if err != nil {
		return "", err
	}
	// remove files from redis after merge
	defer r.remove(ctx)

	// merge files
	finalFile, err := os.CreateTemp(r.tempDir, "*")
	if err != nil {
		return "", err
	}
	defer func() { _ = finalFile.Close() }()

	writer := sha256.New()

	multiWriter := io.MultiWriter(finalFile, writer)

	// merge function is an inner function to merge files
	mergeAndRemove := func(filepath string) error {
		f, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
			_ = os.Remove(filepath)
		}()
		_, err = io.Copy(multiWriter, f)
		return err
	}

	// merge files
	for _, file := range files {
		// merge file and remove it
		if err = mergeAndRemove(file); err != nil {
			return "", err
		}
	}

	// try to check file hash
	// if hash is not equal to fileHash, return error
	if hex.EncodeToString(writer.Sum(nil)) != r.fileHash {
		return "", errors.New("hash not equal")
	}
	return filepath.Base(finalFile.Name()), nil
}

// remove is a method that removes all files from the local filesystem and deletes the Redis list.
func (r *redisLocalFileMerger) remove(ctx context.Context) {
	r.client.Del(ctx, r.fileHash)
}

// do not import "github.com/eatmoreapple/wxhelper/internal/wxclient" directly
//
//go:linkname tempDir github.com/eatmoreapple/wxhelper/internal/wxclient.tempDir
func tempDir() string

// localFileMergerFactory is a struct that holds a Redis client.
type localFileMergerFactory struct {
	client *redis.Client
}

// New is a method that creates a new instance of redisLocalFileMerger.
func (l *localFileMergerFactory) New(key string) (FileMerger, error) {
	items := strings.Split(key, ":")
	if len(items) != 2 {
		return nil, errors.New("invalid key")
	}
	filename, fileHash := items[0], items[1]
	return &redisLocalFileMerger{
		client:   l.client,
		filename: filename,
		fileHash: fileHash,
		tempDir:  tempDir(),
	}, nil
}
