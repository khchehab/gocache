package gocache

import "time"

// cacheValue is a structure that represents the cache value.
// It contains the actual value, the TTL and the expiry date of the value.
type cacheValue struct {
	// value is the actual value of the cache entry.
	value any
	// size is the size of the value in bytes.
	size uint64
	// ttl is the time-to-live duration of the cache value entry.
	ttl time.Duration
	// expiryDate is the cache entry value expiration date.
	expiryDate time.Time
	// timer is timer of the cache value if delete on expire is set on it.
	timer *time.Timer
}

// Expired returns a flag whether the cache entry has expired or not.
//
// Returns:
//   - bool: A flag if the cache entry has expired or not.
func (v *cacheValue) Expired() bool {
	return v.ttl > 0 && v.expiryDate.Before(time.Now().UTC())
}

// CacheEntry is a structure that represents a cache entry.
// This structure is used for the multiple operations.
type CacheEntry struct {
	// Key is the key of the cache entry.
	Key string
	// Value is the value of the cache entry.
	Value any
	// Ttl is the time-to-live of the cache entry.
	Ttl int
}
