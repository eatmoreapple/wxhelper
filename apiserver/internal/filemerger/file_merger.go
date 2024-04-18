package filemerger

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/eatmoreapple/env"
	"github.com/go-redis/redis/v8"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileMerger interface {
	Add(ctx context.Context, file string) error
	Merge() (string, error)
}

type Factory interface {
	New(key string) (FileMerger, error)
}

type redisLocalFileMerger struct {
	client   *redis.Client
	filename string
	fileHash string
}

func (r *redisLocalFileMerger) Add(ctx context.Context, file string) error {
	return r.client.LPush(ctx, r.fileHash, file).Err()
}

func (r *redisLocalFileMerger) Merge() (string, error) {
	// try to get all files from redis
	files, err := r.client.LRange(context.Background(), r.fileHash, 0, -1).Result()
	if err != nil {
		return "", err
	}
	// remove files from redis after merge
	defer r.remove(files)

	// merge files
	finalFile, err := os.CreateTemp(env.Name("TEMP_DIR").String(), "*")
	if err != nil {
		return "", err
	}
	defer func() { _ = finalFile.Close() }()

	writer := sha256.New()

	// merge function is an inner function to merge files
	merge := func(f io.Reader) error {
		if closer, ok := f.(io.Closer); ok {
			defer func() { _ = closer.Close() }()
		}
		_, err := io.Copy(finalFile, io.TeeReader(f, writer))
		return err
	}

	// merge files
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return "", err
		}
		if err = merge(f); err != nil {
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

func (r *redisLocalFileMerger) remove(files []string) {
	for _, file := range files {
		_ = os.Remove(file)
	}
	r.client.Del(context.Background(), r.fileHash)
}

type localFileMergerFactory struct {
	client *redis.Client
}

func (l *localFileMergerFactory) New(key string) (FileMerger, error) {
	items := strings.Split(key, ":")
	if len(items) != 2 {
		return nil, errors.New("invalid key")
	}
	filename, fileHash := items[0], items[1]
	return &redisLocalFileMerger{client: l.client, filename: filename, fileHash: fileHash}, nil
}

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
