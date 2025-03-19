package gocache

import (
	"maps"
	"slices"
	"sync"
	"time"
)

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

	data  map[string]*cacheValue
	mu    sync.Mutex
	stats *Stats
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
		stats: &Stats{
			Hits:      0,
			Misses:    0,
			Keys:      0,
			KeySize:   0,
			ValueSize: 0,
		},
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

	if !ok && c.maxKeys != -1 && int(c.stats.Keys) >= c.maxKeys {
		return ErrCacheFull
	}

	valueSize := SizeOf(value)

	if ok {
		if val.timer != nil {
			val.timer.Stop()
		}

		c.stats.ValueSize -= val.size
		delete(c.data, key)
	} else {
		c.stats.Keys++
		c.stats.KeySize += SizeOf(key)
		c.stats.ValueSize += valueSize
	}

	keyTtl := c.stdTtl
	if ttl > -1 {
		keyTtl = ttl
	}

	expiryDate := time.Now().UTC().Add(keyTtl)
	val = &cacheValue{
		value:      value,
		size:       valueSize,
		ttl:        keyTtl,
		expiryDate: expiryDate,
		timer:      nil,
	}
	c.data[key] = val

	if keyTtl > 0 && c.deleteOnExpire {
		c.data[key].timer = time.AfterFunc(keyTtl, func() {
			c.mu.Lock()
			defer c.mu.Unlock()

			c.stats.Keys--
			c.stats.KeySize -= SizeOf(key)
			c.stats.ValueSize -= val.size
			delete(c.data, key)
		})
	}

	return nil
}

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

	val, ok := c.data[key]

	if !ok {
		c.stats.Misses++
		return nil, ErrKeyNotFound
	}

	if val.Expired() {
		c.stats.Misses++
		return nil, ErrKeyNotFound
	}

	c.stats.Hits++
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

	val, ok := c.data[key]

	if !ok {
		c.stats.Misses++
		return nil, ErrKeyNotFound
	}

	if val.Expired() {
		c.stats.Misses++
		return nil, ErrKeyNotFound
	}

	if val.timer != nil {
		val.timer.Stop()
	}

	c.stats.Hits++
	c.stats.Keys--
	c.stats.KeySize -= SizeOf(key)
	c.stats.ValueSize -= val.size
	delete(c.data, key)

	return val.value, nil
}

// Delete removes the specified key from the cache if it exists.
// It uses a mutex lock to ensure thread-safety.
//
// The function never fails, it returns the number of key that have been deleted.
// If the key was found and deleted, the function will return `1`, otherwise `0`.
//
// Parameters:
//   - key: The key to delete from the cache.
//
// Returns:
//   - int: The number of entries that have been deleted from the cache.
func (c *Cache) Delete(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0

	val, ok := c.data[key]

	if !ok {
		return 0
	}

	if val.timer != nil {
		val.timer.Stop()
	}

	c.stats.Keys--
	c.stats.KeySize -= SizeOf(key)
	c.stats.ValueSize -= val.size
	delete(c.data, key)
	count++

	return count
}

// ChangeTtl changes the TTL of a key. The function returns whether the TTL has been changed or not.
// It uses a mutex lock to ensure thread-safety.
//
// Below are the possible TTL values:
//   - `-1` will delete the key.
//   - `0` will make the key entry not expire.
//   - Any other value will set the TTL for the key entry.
//
// Parameters:
//   - key: The key to change the TTL of.
//   - ttl: The TTL to update.
//
// Returns:
//   - bool: A boolean flag whether the TTL of the key has been changed or not.
func (c *Cache) ChangeTtl(key string, ttl time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]

	if !ok || val.Expired() {
		return false
	}

	if val.timer != nil {
		val.timer.Stop()
	}

	if ttl < 0 {
		c.stats.Keys--
		c.stats.KeySize -= SizeOf(key)
		c.stats.ValueSize -= val.size
		delete(c.data, key)
		return true
	}

	val.ttl = ttl
	val.expiryDate = time.Now().UTC().Add(ttl)

	if ttl > 0 && c.deleteOnExpire {
		val.timer = time.AfterFunc(ttl, func() {
			c.mu.Lock()
			defer c.mu.Unlock()

			c.stats.Keys--
			c.stats.KeySize -= SizeOf(key)
			c.stats.ValueSize -= val.size
			delete(c.data, key)
		})
	}

	return true
}

// GetTtl returns the TTL (time-to-live) for the specified key.
// It uses a mutex lock to ensure thread-safety.
//
// Cases:
//   - If the key does not exist, a value of `-1` is returned.
//   - If the key exists and has no TTL, a value of `0` is returned.
//   - If the key exists and has a TTL, the TTL value is returned.
//
// Parameters:
//   - key: The key to get the TTL of.
//
// Returns:
//   - time.Duration: The TTL of the key, `0` if no TTL, and `-1` if the key does not exist.
func (c *Cache) GetTtl(key string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]

	if !ok || val.Expired() {
		return -1
	}

	return val.ttl
}

// Keys returns a slice of all keys currently stored in the cache.
// It uses a mutex lock to ensure thread-safety.
//
// Returns:
//   - []string: A slice containing all keys in the cache.
func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

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
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.data[key]
	return ok
}

// Stats returns a copy of the current cache statistics such as hits, misses, key count, the total key size and the total value size.
//
// Returns:
//   - Stats: A copy of the cache statistics.
func (c *Cache) Stats() Stats {
	return *c.stats
}

// Clear removes all key-value entries from the cache and resets statistics.
//
// This function safely clears the entire cache by:
//   - Stopping any active timers associated with expiring keys.
//   - Deleting all entries from the cache.
//   - Resetting cache statistics to zero.
//
// It uses a mutex lock to ensure thread-safety.
//
// Usage:
//   - Call this function to completely reset the cache, removing all stored data.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.data {
		if v.timer != nil {
			v.timer.Stop()
		}

		delete(c.data, k)
	}

	c.stats = &Stats{
		Hits:      0,
		Misses:    0,
		Keys:      0,
		KeySize:   0,
		ValueSize: 0,
	}
}

// ClearStats resets all cache statistics to zero in a thread-safe manner.
//
// This function clears the stored statistics, including hits, misses, key count,
// key size, and value size. It does not affect the actual cached data.
//
// It uses a mutex lock to ensure thread-safety.
//
// Usage:
//   - Call this function to reset cache statistics, e.g., after a performance measurement.
func (c *Cache) ClearStats() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats = &Stats{
		Hits:      0,
		Misses:    0,
		Keys:      0,
		KeySize:   0,
		ValueSize: 0,
	}
}
