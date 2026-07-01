package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	cache := NewCache(1)

	// Test adding items to the cache
	cache.Add("key1", []byte("value1"))
	cache.Add("key2", []byte("value2"))
	cache.Add("key3", []byte("value3"))

	// Test retrieving items from the cache
	if val, ok := cache.Get("key1"); !ok || string(val) != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	if val, ok := cache.Get("key2"); !ok || string(val) != "value2" {
		t.Errorf("Expected value2, got %v", val)
	}

	if val, ok := cache.Get("key3"); !ok || string(val) != "value3" {
		t.Errorf("Expected value3, got %v", val)
	}

	time.Sleep(2 * time.Second) // Wait for the cache to expire

	// Test that items have been removed from the cache
	if _, ok := cache.Get("key1"); ok {
		t.Errorf("Expected key1 to be removed from cache")
	}

	if _, ok := cache.Get("key2"); ok {
		t.Errorf("Expected key2 to be removed from cache")
	}

	if _, ok := cache.Get("key3"); ok {
		t.Errorf("Expected key3 to be removed from cache")
	}

}
