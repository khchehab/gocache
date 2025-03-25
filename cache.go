package gocache

import (
	"time"
)

// Cache is an in-memory key-value store.
// The cache contains configurations that dictate its behavior, below are the default values:
//   - StdTTL: 0 - entries never expire.
//   - DeleteOnExpire: true - entries are automatically deleted upon expiration.
//   - MaxKeys: -1 - unlimited number of entries.
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
}

// New creates a new [Cache] instance with optional configurations and an empty data store.
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

// Set sets a key-value pair in the cache.
// If an error occurs, it will be returned, otherwise nil will be returned.
func (c *Cache) Set(key string, value any) error {
	return c.SetWithTtl(key, value, -1)
}

// SetWithTtl sets a key-value pair in the cache with a TTL (time-to-live) in duration.
// If an error occurs, it will be returned, otherwise nil will be returned.
func (c *Cache) SetWithTtl(key string, value any, ttl time.Duration) error {
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
	val = &cacheValue{
		value:      value,
		ttl:        keyTtl,
		expiryDate: expiryDate,
		timer:      nil,
	}
	c.data[key] = val

	if keyTtl > 0 && c.deleteOnExpire {
		c.data[key].timer = time.AfterFunc(keyTtl, func() {
			delete(c.data, key)
		})
	}

	return nil
}

// Get returns the value associated with the provided key from the cache.
// It returns the value if found in the cache.
// If an error occurs, it will be returned, otherwise nil will be returned.
func (c *Cache) Get(key string) (any, error) {
	val, ok := c.data[key]

	if !ok {
		return nil, ErrKeyNotFound
	}

	if val.expired() {
		return nil, ErrKeyNotFound
	}

	return val.value, nil
}

// GetAndDelete returns the value associated with the provided key from the cache and removes it.
// It returns the value if found in the cache.
// If an error occurs, it will be returned, otherwise nil will be returned.
func (c *Cache) GetAndDelete(key string) (any, error) {
	val, ok := c.data[key]

	if !ok {
		return nil, ErrKeyNotFound
	}

	if val.expired() {
		return nil, ErrKeyNotFound
	}

	if val.timer != nil {
		val.timer.Stop()
	}

	delete(c.data, key)

	return val.value, nil
}

// Delete removes the entry associated with the provided key from the cache if it exists.
// It returns the number of deleted items from the cache.
func (c *Cache) Delete(key string) int {
	count := 0

	val, ok := c.data[key]

	if !ok {
		return 0
	}

	if val.timer != nil {
		val.timer.Stop()
	}

	delete(c.data, key)
	count++

	return count
}

// ChangeTtl changes the TTL associated with the provided key in the cache.
// It returns a bool indicating whether a change in TTL has occurred or not.
func (c *Cache) ChangeTtl(key string, ttl time.Duration) bool {
	val, ok := c.data[key]

	if !ok || val.expired() {
		return false
	}

	if val.timer != nil {
		val.timer.Stop()
	}

	if ttl < 0 {
		delete(c.data, key)

		return true
	}

	val.ttl = ttl
	val.expiryDate = time.Now().UTC().Add(ttl)

	if ttl > 0 && c.deleteOnExpire {
		val.timer = time.AfterFunc(ttl, func() {
			delete(c.data, key)
		})
	}

	return true
}

// GetTtl returns the TTL, as a duration, of the provided key in the cache.
// It returns -1 if the key does not exist.
func (c *Cache) GetTtl(key string) time.Duration {
	val, ok := c.data[key]

	if !ok || val.expired() {
		return -1
	}

	return val.ttl
}

// Keys returns the list of keys, as a slice of string, in the cache.
func (c *Cache) Keys() []string {
	keys := make([]string, 0, len(c.data))

	for k := range c.data {
		keys = append(keys, k)
	}

	return keys
}

// Has returns a bool whether the key exists in the cache or not.
func (c *Cache) Has(key string) bool {
	val, ok := c.data[key]

	if !ok || val.expired() {
		return false
	}

	return true
}

// Clear clears the cache by emptying the store.
func (c *Cache) Clear() {
	if len(c.data) == 0 {
		return
	}

	for _, v := range c.data {
		if v.timer != nil {
			v.timer.Stop()
		}
	}

	c.data = make(map[string]*cacheValue)
}
