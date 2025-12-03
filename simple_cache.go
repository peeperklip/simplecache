package keyvalstore

import (
	"sync"
	"time"
)

// SimpleCache is a thread-safe in-memory key-value store with expiration.
type SimpleCache[T any] struct {
	data      map[string]cacheItem[T]
	mutex     sync.RWMutex
	expiryDur time.Duration
}

type cacheItem[T any] struct {
	value      T
	expiryTime time.Time
}

// NewSimpleCache creates a new SimpleCache with the specified expiration duration.
func NewSimpleCache[T any](expiryDur time.Duration) *SimpleCache[T] {
	return &SimpleCache[T]{
		data:      make(map[string]cacheItem[T]),
		expiryDur: expiryDur,
	}
}

// Set adds a key-value pair to the cache with an expiration time.
func (c *SimpleCache[T]) Set(key string, value T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = cacheItem[T]{
		value:      value,
		expiryTime: time.Now().Add(c.expiryDur),
	}
}

// Get retrieves a value from the cache by key.
// It returns the value and a boolean indicating whether the key was found and not expired.
func (c *SimpleCache[T]) Get(key string) (T, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	item, exists := c.data[key]
	var zero T
	if !exists {
		return zero, false
	}

	if time.Now().After(item.expiryTime) {
		delete(c.data, key)
		return zero, false
	}

	return item.value, true
}
