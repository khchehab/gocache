package gocache

import "time"

// OptFunc defines a function type for configuring a Cache instance.
type OptFunc func(*Cache)

// WithStdTtl returns an OptFunc that sets the cache's global TTL (time-to-live).
// If TTL is reached for each entry (based on it's own insertion time), it will indicate the entry has expired.
// A value of `0` means unlimited.
//
// Parmeters:
//   - stdTtl: The TTL for all cache entries.
//
// Returns:
//   - OptFunc: A function that applies the stdTtl setting to a Cache instance.
func WithStdTtl(stdTtl time.Duration) OptFunc {
	return func(c *Cache) {
		if stdTtl > -1 {
			c.stdTtl = stdTtl
		}
	}
}

// WithDeleteOnExpire returns an OptFunc that sets the flag whether to automatically
// delete expires keys or keep them flagged only.
//
// Parameters:
//   - deleteOnExpire: A boolean value to indicate whether to delete key when expired or not.
//
// Returns:
//   - OptFunc: A function thay applies the deleteOnExpire setting to a Cache instance.
func WithDeleteOnExpire(deleteOnExpire bool) OptFunc {
	return func(c *Cache) {
		c.deleteOnExpire = deleteOnExpire
	}
}

// WithMaxKeys returns an OptFunc that sets the maximum number of keys allowed in the cache.
// If maxKeys is set ot a positive number, the cache will enforce this limit.
// A negative value (e.g. -1) indicates no limit.
//
// Parameters:
//   - maxKeys: The maximum number of keys the cache can store.
//
// Returns:
//   - OptFunc: A function that applies the maxKeys setting to a Cache instance.
func WithMaxKeys(maxKeys int) OptFunc {
	return func(c *Cache) {
		if maxKeys > -1 {
			c.maxKeys = maxKeys
		}
	}
}
