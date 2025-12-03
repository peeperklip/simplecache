package keyvalstore

import (
	"sync"
	"testing"
	"time"
)

func TestSimpleCache_SetAndGet(t *testing.T) {
	sut := NewSimpleCache[string](1 * time.Second)

	sut.Set("key1", "value1")
	val, found := sut.Get("key1")
	if !found || val != "value1" {
		t.Errorf("Expected to find key1 with value 'value1', got '%s', found: %v", val, found)
	}

	val, found = sut.Get("key2")
	if found {
		t.Errorf("Expected not to find key2, but got value '%s'", val)
	}
}

func TestSimpleCache_DurationSetToZeroWillNotCache(t *testing.T) {
	sut := NewSimpleCache[string](0)

	sut.Set("key1", "value1")
	val, found := sut.Get("key1")
	if found {
		t.Errorf("Expected not to find key1, but got value '%s'", val)
	}
}

func TestSimpleCache_Expiration(t *testing.T) {
	sut := NewSimpleCache[string](10 * time.Millisecond)

	sut.Set("key1", "value1")
	val, found := sut.Get("key1")
	if !found || val != "value1" {
		t.Errorf("Expected to find key1 with value 'value1', got '%s', found: %v", val, found)
	}

	time.Sleep(15 * time.Millisecond)
	val, found = sut.Get("key1")
	if found {
		t.Errorf("Expected key1 to be expired, but got value '%s'", val)
	}
}

func TestSimpleCache_ConcurrentAccess(t *testing.T) {
	sut := NewSimpleCache[int](10 * time.Second)
	const numGoroutines = 100
	const numIterations = 1000

	var wg sync.WaitGroup
	wg.Add(2 * numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				sut.Set("key", id*j)
			}
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				sut.Get("key")
			}
		}()
	}

	wg.Wait()
}

func TestSimpleCache_DifferentTypes(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "string cache",
			run: func(t *testing.T) {
				c := NewSimpleCache[string](1 * time.Minute)
				c.Set("strKey", "stringValue")
				v, ok := c.Get("strKey")
				if !ok || v != "stringValue" {
					t.Errorf("Expected to find strKey with value 'stringValue', got '%s', found: %v", v, ok)
				}
			},
		},
		{
			name: "int cache",
			run: func(t *testing.T) {
				c := NewSimpleCache[int](1 * time.Minute)
				c.Set("intKey", 42)
				v, ok := c.Get("intKey")
				if !ok || v != 42 {
					t.Errorf("Expected to find intKey with value 42, got '%d', found: %v", v, ok)
				}
			},
		},
		{
			name: "struct cache",
			run: func(t *testing.T) {
				type testStruct struct {
					Field1 string
					Field2 int
				}
				c := NewSimpleCache[testStruct](1 * time.Minute)
				expected := testStruct{Field1: "test", Field2: 100}
				c.Set("structKey", expected)
				v, ok := c.Get("structKey")
				if !ok || v != expected {
					t.Errorf("Expected to find structKey with value %+v, got %+v, found: %v", expected, v, ok)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.run)
	}
}
