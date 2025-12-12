package keyvalstore

import (
	"sync"
	"time"
)

// SimpleCache is a thread-safe in-memory key-value store with expiration.
type SimpleCache[T any] struct {
	data            map[string]cacheItem[T]
	cleanupInterval time.Duration

	mutex     sync.RWMutex
	done      chan struct{}
	wg        sync.WaitGroup
	closeOnce sync.Once
}

type cacheItem[T any] struct {
	value      T
	expiryTime time.Time
}

// NewSimpleCache creates a new SimpleCache with a specified cleanup interval.
func NewSimpleCache[T any](cleanupInterval time.Duration) *SimpleCache[T] {
	c := &SimpleCache[T]{
		data:            make(map[string]cacheItem[T]),
		done:            make(chan struct{}),
		cleanupInterval: cleanupInterval,
	}

	c.wg.Add(1)
	go c.janitor()
	return c
}

// Set adds a key-value pair to the cache with an expiration time.
func (c *SimpleCache[T]) Set(key string, expiryDur time.Duration, value T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = cacheItem[T]{
		value:      value,
		expiryTime: time.Now().Add(expiryDur),
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
		// Item has expired, return zero value and false.
		// Note: Expired items will be cleaned up by the janitor goroutine.
		return zero, false
	}

	return item.value, true
}

func (c *SimpleCache[T]) janitor() {
	defer c.wg.Done()
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			c.mutex.Lock()
			for k, it := range c.data {
				if now.After(it.expiryTime) {
					delete(c.data, k)
				}
			}
			c.mutex.Unlock()
		case <-c.done:
			return
		}
	}
}

// Close stops the janitor goroutine and waits for it to exit.
func (c *SimpleCache[T]) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
	})
	c.wg.Wait()
}
