package gocache

import "time"

// OptFunc defines a function type for configuring a [Cache] instance.
type OptFunc func(*Cache)

// WithStdTtl returns an [OptFunc] that sets the global cache's TTL (time-to-live).
// It takes a duration as a parameter and updates the cache accordingly.
// A duration of 0 means unlimited, the keys never expire.
func WithStdTtl(stdTtl time.Duration) OptFunc {
	return func(c *Cache) {
		if stdTtl > -1 {
			c.stdTtl = stdTtl
		}
	}
}

// WithDeleteOnExpire returns an [OptFunc] that sets the cache's flag of deletion on expiry.
// If set to true, entries will be automatically when they expire.
// If set to false, entries will remain in the store but flagged as expired.
func WithDeleteOnExpire(deleteOnExpire bool) OptFunc {
	return func(c *Cache) {
		c.deleteOnExpire = deleteOnExpire
	}
}

// WithMaxKeys returns an [OptFunc] that sets the cache's maximum number of keys.
// A value of -1 means unlimited keys.
func WithMaxKeys(maxKeys int) OptFunc {
	return func(c *Cache) {
		if maxKeys > -1 {
			c.maxKeys = maxKeys
		}
	}
}
