# SimpleCache 
A small, thread-safe in-memory key-value cache with per-entry expiration and Go generics.

### Features
- Generic cache: `SimpleCache[T any]`
- Safe for concurrent use (uses `sync.RWMutex`)
- Per-entry expiration based on a common `expiryDur`
- Simple API: `NewSimpleCache`, `Set`, `Get`

### Limitations
- No eviction policy beyond expiration
- Not persistent; data is lost on program exit

### Usage

```go
package main

import (
	"fmt"
	"time"

	`github.com/peeperklip/simplecache`
)

func main() {
	var (
		cleanupInterval = time.Minute * 5
		greetingTTL = time.Minute * 10
	)
	// Create a cache that cleans up expired entries every 5 minutes
	cache := keyvalstore.NewSimpleCache[string](cleanupInterval)

	// Store a value with a TTL
	cache.Set("greeting", greetingTTL, "hello")

	// Retrieve a value
	if v, ok := cache.Get("greeting"); ok {
		fmt.Println("value:", v)
	} else {
		fmt.Println("missing or expired")
	}
}
```