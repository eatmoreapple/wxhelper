// Package filemerger provides interfaces and functions for merging files.
package filemerger

import (
	"context"
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
