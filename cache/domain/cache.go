package domain

import (
	"fmt"
	"sync"
	"time"
)

const DefaultTTLExpiration = 30

type Cache interface {
	Set(key string, value string) error
	Get(key string) (*Item, error)
}

type Item struct {
	Value      string
	Expiration time.Time
}

type SimpleCache struct {
	lock sync.RWMutex
	data map[string]Item
}

func NewCache() *SimpleCache {
	cache := &SimpleCache{
		data: make(map[string]Item),
	}
	go cache.cleanCache()
	return cache
}

func (sc *SimpleCache) cleanCache() {
	for {
		time.Sleep(DefaultTTLExpiration * time.Minute)
		sc.lock.Lock()
		for key, item := range sc.data {
			if time.Now().After(item.Expiration) {
				delete(sc.data, key)
			}
		}
		sc.lock.Unlock()
	}
}

func (sc *SimpleCache) Get(key string) (*Item, error) {
	sc.lock.RLock()
	defer sc.lock.RUnlock()
	item, ok := sc.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}
	return &item, nil
}

func (sc *SimpleCache) Set(key string, value string) error {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	sc.data[key] = Item{value, time.Now().Add(DefaultTTLExpiration * time.Minute)}
	return nil
}
