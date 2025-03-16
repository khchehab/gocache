package gocache

import "errors"

// ErrKeyNotFound is an error for when a key doesn't exist in the cache.
var ErrKeyNotFound = errors.New("key not found")

// ErrCacheFull is an error for when the cache has reached the maximum allowed number of items.
var ErrCacheFull = errors.New("the cache is full")
