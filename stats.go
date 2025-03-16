package gocache

// Stats represents cache performance metrics and storage statistics.
// It tracks the number of cache hits, misses, and key/value storage details.
type Stats struct {
	// Hits is the number of times a requested key was found in the cache.
	Hits uint
	// Misses is the number of times a requested key was not found in the cache.
	Misses uint
	// Keys is the total number of keys currently stored in the cache.
	Keys uint
	// KeySize is the total size (in bytes) of all keys stored in the cache.
	KeySize uint
	// ValueSize is the total size (in bytes) of all values stored in the cache.
	ValueSize uint
}
