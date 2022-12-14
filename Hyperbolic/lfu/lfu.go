package cache

import (
	"math"
)

// A CacheItem is an item that will go in the cache
type LFUCacheItem struct {

	// size of value (bytes)
	value_size int

	// how many times the item is accessed
	access_count int

}

// LFU is a cache that follows the LFU algorithm
type LFU struct {

	// max number of bytes the cache can hold
	max_capacity int

	// total size of items that are stored in the cache
	size int

	// number of bindings in cache
	num_bindings int

	// map from keys to items in cache
	mapping map[string]*LFUCacheItem

	// number of hits
	hits int

	// number of misses
	misses int
}

// NewHyperbolicCache creates a new, empty Hyperbolic cache
func newLFU(max_capacity int) *LFU {

	return &LFU{
		max_capacity: max_capacity,
		size:         0,
		num_bindings: 0,
		mapping:      make(map[string]*LFUCacheItem, max_capacity),
		hits:         0,
		misses:       0,
	}
}

// Given a key, Get return the corresponding value's size and a success boolean
func (cache *LFU) Get(key string) (value_size int, ok bool) {

	// retrieve item associated with key
	item, ok := cache.mapping[key]

	value_size = 0

	if ok {
		cache.hits += 1

		// update access count of item
		item.access_count += 1

		// return size of value
		value_size = item.value_size

	} else {
		cache.misses += 1

		return 0, false
	}

	return value_size, ok
}

// Given a key and value, Set add them to the cache
func (cache *LFU) Set(key string, value_size int) bool {

	// size of value to be added
	value_length := value_size

	// size of key and value to be added
	insert_size := (len([]byte(key)) + value_length)

	// key and value pair is too large
	if (insert_size) > cache.max_capacity {
		return false
	}

	// check if the key already has a value
	existing_item, ok := cache.mapping[key]
	if ok {
		// replace value size, update size of cache, and return
		cache.size -= existing_item.value_size
		cache.size += value_length

		// update value size and access count of the actual item
		cache.mapping[key].value_size = value_length
		cache.mapping[key].access_count += 1

		// remove until the size meets the capacity 
		for (cache.size > cache.max_capacity) {
			key_to_remove := cache.Evict_Which()
			cache.size -= (cache.mapping[key_to_remove].value_size + len([]byte(key_to_remove)))
			_, success := cache.Remove(key_to_remove)

			if success {
				cache.num_bindings -= 1
			}
		}

		return true
	}

	// if not enough space, evict until there is enough space
	for insert_size > (cache.max_capacity - cache.size) {

		// find what key, value pair we should evict
		key_to_remove := cache.Evict_Which()

		// remove the chosen key
		_, success := cache.Remove(key_to_remove)
		if success {
			cache.num_bindings -= 1
		}
	}

	cache.mapping[key] = &LFUCacheItem{value_size: value_length, access_count: 1,}

	cache.size += insert_size
	cache.num_bindings += 1

	return true
}


// Evict_Which() is an algorithm to select which item in the cache to evict
func (cache *LFU) Evict_Which() (key string) {

	if len(cache.mapping) == 0 {
		return ""
	}

	// iterate through mapping to find the item with 
	// the least number of accesses
	minimum := ""
	minValue := math.MaxInt32

	for j := range cache.mapping {
		if cache.mapping[j].access_count < minValue {
			minValue = cache.mapping[j].access_count
			minimum = j
		}
	}

	return minimum
}

// MaxStorage returns the maximum number of bytes this cache can store
func (cache *LFU) MaxStorage() int {
	return cache.max_capacity
}

// RemainingStorage returns the number of unused bytes available in this cache
func (cache *LFU) RemainingStorage() int {
	return cache.max_capacity - cache.size
}

// Stats returns statistics about how many search hits and misses have occurred.
func (cache *LFU) Stats() *Stats {
	return &Stats{Hits: cache.hits, Misses: cache.misses}
}

// Len returns the number of bindings in the cache.
func (cache *LFU) Len() int {
	return cache.num_bindings
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value size was found and false otherwise
func (cache *LFU) Remove(key string) (value_size int, ok bool) {

	// check if there is an item associated with key
	item, ok := cache.mapping[key]
	if !ok {
		return 0, false
	}

	// update the capacity
	cache.size -= (cache.mapping[key].value_size + len([]byte(key)))

	// remove from hashmap
	delete(cache.mapping, key)

	return item.value_size, true
}
