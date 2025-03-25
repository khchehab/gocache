package gocache

import "time"

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

// expired returns a flag whether the cache entry has expired or not.
func (v *cacheValue) expired() bool {
	return v.ttl > 0 && v.expiryDate.Before(time.Now().UTC())
}
