package cache

import (
	"log"
	"math"
)

// A LFUCacheItem is an item with metadata that implicitly
// holds a value of set size. It goes in a LFUCache.
type LFUCacheItem struct {

	// how many times the item has been accessed
	access_count int
}

// A LFUCache is a cache that uses lfu caching.
type LFUCache struct {

	// maximum number of items the cache can hold
	max_capacity int

	// total number of items currently in the cache
	size int

	// map of keys to items in the cache
	keys_to_items map[string]*LFUCacheItem

	// number of hits
	hits int

	// number of misses
	misses int
}

// NewLFUCache creates a new, empty LFUCache.
func NewLFUCache(max_capacity int) *LFUCache {

	return &LFUCache{
		max_capacity:  max_capacity,
		size:          0,
		keys_to_items: make(map[string]*LFUCacheItem, max_capacity),
		hits:          0,
		misses:        0,
	}
}

// Get returns how many times the item associated with the given key
// has been accessed (not including this access) and a success boolean.
func (cache *LFUCache) Get(key string) (access_count int, ok bool) {

	// retrieve item associated with key
	item, ok := cache.keys_to_items[key]

	if ok {
		cache.hits += 1

		access_count = item.access_count

		// update access count of item
		item.access_count += 1

	} else {
		cache.misses += 1

		return 0, false
	}

	return access_count, ok
}

// Set adds/updates an item with the given key in the cache
// and returns a success boolean.
func (cache *LFUCache) Set(key string) (ok bool) {

	// check if an item with that key already exists
	existing_item, ok := cache.keys_to_items[key]

	if ok {
		// update access count of item
		existing_item.access_count += 1

		return true
	}

	// if not enough space and an item with the key does not exist,
	// evict an item
	if cache.size == cache.max_capacity {

		key_to_remove := cache.Evict_Which()

		success := cache.Remove(key_to_remove)
		if !success {
			log.Fatal("Failed to evict an item.")
		}

	}

	// add new item with key
	cache.keys_to_items[key] = &LFUCacheItem{
		access_count: 1}

	// update size of cache
	cache.size += 1

	return true
}

// Evict_Which() is an algorithm to select which item in the cache to evict.
func (cache *LFUCache) Evict_Which() (key string) {

	// make sure cache is actually full before evicting
	if cache.size != cache.max_capacity {
		log.Fatal("Should not be evicting when cache is not full.")
	}

	// iterate through keys_to_items to find the item with
	// the least number of accesses
	minimum := ""
	minValue := math.MaxInt32

	for j := range cache.keys_to_items {
		if cache.keys_to_items[j].access_count < minValue {
			minValue = cache.keys_to_items[j].access_count
			minimum = j
		}
	}

	return minimum
}

// MaxStorage returns the maximum number of items this cache can store.
func (cache *LFUCache) MaxStorage() int {
	return cache.max_capacity
}

// RemainingStorage returns the number of items that can still be stored
// in this cache.
func (cache *LFUCache) RemainingStorage() int {
	return cache.max_capacity - cache.size
}

// Stats returns statistics about how many search hits and misses have occurred.
func (cache *LFUCache) Stats() *Stats {
	return &Stats{Hits: cache.hits, Misses: cache.misses}
}

// Len returns the number of items in the cache.
func (cache *LFUCache) Len() int {
	return cache.size
}

// Remove removes the item associated with the given key from the cache, if it exists.
// ok is true if an item was found and false otherwise.
func (cache *LFUCache) Remove(key string) (ok bool) {

	// check if there is an item associated with key
	_, ok = cache.keys_to_items[key]
	if !ok {
		return false
	}

	// remove the key from the cache
	delete(cache.keys_to_items, key)

	cache.size -= 1

	return true
}
