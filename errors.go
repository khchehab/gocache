package gocache

import "errors"

var (
	// ErrKeyNotFound is an error for when a key doesn't exist in the cache.
	ErrKeyNotFound = errors.New("key not found")

	// ErrCacheFull is an error for when the cache has reached the maximum allowed number of items.
	ErrCacheFull = errors.New("the cache is full")
)
