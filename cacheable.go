package gocache

import "time"

type Cacheable interface {
	Set(key string, value any) error
	SetWithTtl(key string, value any, ttl time.Duration) error
	Get(key string) (any, error)
	GetAndDelete(key string) (any, error)
	Delete(key string) int
	ChangeTtl(key string, ttl time.Duration) bool
	GetTtl(key string) time.Duration
	Keys() []string
	Has(key string) bool
	Clear()
}
