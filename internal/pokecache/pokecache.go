package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entry map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

func NewCache(interval int) *Cache {
	in := time.Duration(interval) * time.Second
	c := &Cache{
		entry: make(map[string]cacheEntry),
	}
	c.readLoop(in)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.entry[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) readLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	go func() {
		for range t.C {
			c.mu.Lock()
			for key, entry := range c.entry {
				if time.Since(entry.createdAt) > interval {
					delete(c.entry, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}
