package filemerger

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/eatmoreapple/wxhelper/apiserver/internal/filemerger/internal"
	"io"
	"os"
	"path/filepath"
	_ "unsafe"
)

// localFileMerger is a struct that holds a Redis client, filename and fileHash.
type localFileMerger struct {
	cache    internal.Cache
	filename string
	fileHash string
	tempDir  string
}

// Add is a method that adds a file to a Redis list.
func (r *localFileMerger) Add(ctx context.Context, file string) error {
	return r.cache.Add(ctx, r.fileHash, file)
}

// Merge is a method that merges all files in a Redis list and checks if the hash of the merged file matches the fileHash.
func (r *localFileMerger) Merge(ctx context.Context) (string, error) {
	// try to get all files from redis
	files, err := r.cache.GetAll(ctx, r.fileHash)
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
func (r *localFileMerger) remove(ctx context.Context) {
	_ = r.cache.DelAll(ctx, r.fileHash)
}
