package gocache

import (
	"log"
	"maps"
	"slices"
	"sync"
	"time"
)

// cacheValue is a structure that represents the cache value.
// It contains the actual value, the TTL and the expiry date of the value.
type cacheValue struct {
	// value is the actual value of the cache entry.
	value any
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

// Cache is a thread-safe in-memory key-value store.
type Cache struct {
	// StdTtl defines the time-to-live for all the cache entries.
	// The value `0` means unlimited.
	stdTtl time.Duration
	// deleteOnExpire defines whether the key should be automatically deleted when
	// it expires or just flagged that its expired.
	deleteOnExpire bool
	// maxKeys defines the maximum number of keys the cache can store.
	// If the cache exceeds this limit, an error will be thrown.
	// The value `-1` means unlimited.
	maxKeys int

	data map[string]*cacheValue
	mu   sync.RWMutex
}

// New creates and returns a new Cache instance with optional configuration and an empty data store.
//
// It accepts a variable number of option functions `OptFunc` which allows customizing the cache.
//
// Parameters:
//   - opts: A variadic list of option functions that modify the cache instance.
//
// Returns:
//   - *Cache: A pointer to the newly created Cache instance.
func New(opts ...OptFunc) *Cache {
	c := &Cache{
		stdTtl:         0,
		deleteOnExpire: true,
		maxKeys:        -1,
		data:           make(map[string]*cacheValue),
	}

	for _, fn := range opts {
		fn(c)
	}

	return c
}

// Set inserts or updates a key-value pair in the cache.
// It uses a mutex lock to ensure thread-safety.
//
// Parameters:
//   - key: The key to store the value under.
//   - value: The value to be stored in the cache.
//
// Returns:
//   - error: `ErrCacheFull` if the cache has reached the maximum number of keys allowed, otherwise `nil`.
func (c *Cache) Set(key string, value any) error {
	return c.SetWithTtl(key, value, -1)
}

// SetWithTtl inserts or updates a key-value pair in the cache with the provided TTL for that entry.
// It uses a mutex lock to ensure thread-safety.
//
// Parameters:
//   - key: The key to store the value under.
//   - value: The value to be stored in the cache.
//   - ttl: The time-to-live for the inserted key-value pair.
//
// Returns:
//   - error: `ErrCacheFull` if the cache has reached the maximum number of keys allowed, otherwise `nil`.
func (c *Cache) SetWithTtl(key string, value any, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]

	if !ok && c.maxKeys != -1 && len(c.data) >= c.maxKeys {
		return ErrCacheFull
	}

	if ok {
		if val.timer != nil {
			val.timer.Stop()
		}

		delete(c.data, key)
	}

	keyTtl := c.stdTtl
	if ttl > -1 {
		keyTtl = ttl
	}

	expiryDate := time.Now().UTC().Add(keyTtl)
	c.data[key] = &cacheValue{
		value:      value,
		ttl:        keyTtl,
		expiryDate: expiryDate,
		timer:      nil,
	}

	if keyTtl > 0 && c.deleteOnExpire {
		c.data[key].timer = time.AfterFunc(keyTtl, func() {
			log.Println("set with ttl - deletion after func - start")
			c.mu.Lock()
			defer c.mu.Unlock()

			log.Println("set with ttl - deletion after func - before obtained lock")

			delete(c.data, key)

			log.Println("set with ttl - deletion after func - deleted entry")
		})
	}

	return nil
}

// TODO MSET (Multiple Set) (Array<(key, value, [ttl])>)

// Get retrieves the value associated with the provided key from the cache in a thread-safe manner.
// It uses a mutex lock to ensure thread-safety.
//
// If the key exists, it returns the value and a nil error.
// If the key is not found, it return nil and ErrKeyNotFound.
//
// Parameters:
//   - key: The key to look up in the cache.
//
// Returns:
//   - any: The value stored in the cache for the provided key.
//   - error: `ErrKeyNotFound` if the key does not exist.
func (c *Cache) Get(key string) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var val *cacheValue
	var ok bool

	if val, ok = c.data[key]; !ok {
		return nil, ErrKeyNotFound
	}

	if val.Expired() {
		if c.deleteOnExpire {
			if val.timer != nil {
				val.timer.Stop()
			}

			delete(c.data, key)
		}

		return nil, ErrKeyNotFound
	}

	return val.value, nil
}

// GetAndDelete retrieves the value associated with the provided key from the cache and removes it.
// It uses a mutex lock to ensure thread-safety.
//
// If the key exists, it returns the value and deletes the entry from the cache.
// If the key is not found, it return nil and ErrKeyNotFound.
//
// Parameters:
//   - key: The key to look up and remove from the cache.
//
// Returns:
//   - any: The value stored in the cache for the given key before deletion.
//   - error: `ErrKeyNotFound` if the key does not exist.
func (c *Cache) GetAndDelete(key string) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var val *cacheValue
	var ok bool

	if val, ok = c.data[key]; !ok {
		return nil, ErrKeyNotFound
	}

	if val.timer != nil {
		val.timer.Stop()
	}

	delete(c.data, key)

	return val.value, nil
}

// Delete removes the specified key(s) from the cache if they exist.
// It uses a mutex lock to ensure thread-safety.
//
// The function never fails, it returns the number of key(s) that have been deleted.
//
// Parameters:
//   - keys: A variadic list of strings for the keys to delete.
//
// Returns:
//   - int: The number of entries that have been deleted from the cache.
func (c *Cache) Delete(keys ...string) int {
	if len(keys) == 0 {
		return 0
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0

	for _, key := range keys {
		if val, ok := c.data[key]; ok {
			if val.timer != nil {
				val.timer.Stop()
			}

			delete(c.data, key)
			count++
		}
	}

	return count
}

// TODO TTL (Change TTL) ([key], ttl) (delete key if ttl < 0)

// TODO getTTL (Get TTL) ([key])

// Keys returns a slice of all keys currently stored in the cache.
// It uses a mutex lock to ensure thread-safety.
//
// Returns:
//   - []string: A slice containing all keys in the cache.
func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return slices.Sorted(maps.Keys(c.data))
}

// Has checks whether the provided key exists in the cache.
// It uses a mutex lock to ensure thread-safety.
//
// Parameters:
//   - key: The key to check its existence.
//
// Returns:
//   - bool: A boolean flag that indicates the existence of the key in the cache.
func (c *Cache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[key]
	return ok
}

// TODO STATS ??
// TODO FLUSH ??
// TODO FLUSH STATS ??
// TODO close cache ??

// TODO event emitters ??? does it work in go?
// events:
// * set
// * del
// * expired
// * flush
// * flush_stats
