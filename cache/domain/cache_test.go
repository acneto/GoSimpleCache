package domain

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheThreadSafety(t *testing.T) {
	cache := NewCache()
	wg := sync.WaitGroup{}

	numThreads := 10
	for i := 0; i < numThreads; i++ {
		wg.Add(1)

		// simulate multiple concurrent requests to the cache
		go func(index int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", index)
			value := fmt.Sprintf("value_%d", index)

			cache.Set(key, value)
			item, _ := cache.Get(key)

			if item.Value != value {
				t.Errorf("Expected %s, got %s", value, item)
			}
		}(i)
	}
	wg.Wait()
}

func TestSimpleCache(t *testing.T) {
	cache := NewCache()

	tests := []struct {
		name      string
		key       string
		item      Item
		expectErr bool
	}{
		{
			name:      "Set and Get a value",
			key:       "key1",
			item:      Item{Value: "value1", Expiration: time.Now().Add(5 * time.Second)},
			expectErr: false,
		},
		{
			name:      "Get non-existent key",
			key:       "nonexistent",
			item:      Item{Value: "value2", Expiration: time.Now().Add(5 * time.Second)},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectErr {
				// Test Set
				err := cache.Set(tt.key, tt.item.Value)
				if err != nil {
					t.Fatalf("Set() error = %v, expectErr %v", err, tt.expectErr)
				}
				// Test Get
				got, err := cache.Get(tt.key)
				if (err != nil) != tt.expectErr {
					t.Fatalf("Get() error = %v, expectErr %v", err, tt.expectErr)
				}
				if !tt.expectErr && got.Value != tt.item.Value {
					t.Errorf("Get() = %v, want %v", got, tt.item.Value)
				}
			} else {
				// Test Get for non-existent key
				_, err := cache.Get(tt.key)
				if (err != nil) != tt.expectErr {
					t.Fatalf("Get() error = %v, expectErr %v", err, tt.expectErr)
				}
			}
		})
	}
}
