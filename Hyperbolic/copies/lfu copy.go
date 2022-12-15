package cache

import (
	"math"
)

// A LFUCacheItem is an item with metadata that
// goes in a LFU.
type LFUCacheItem struct {

	// the size of the item's value (bytes)
	value_size int

	// how many times the item has been accessed
	access_count int
}

// A LFU is a cache that follows the LFU algorithm.
type LFU struct {

	// maximum number of bytes the cache can hold
	max_capacity int

	// total number of bytes currently in the cache
	size int

	// number of bindings in cache
	num_bindings int

	// map of keys to items in the cache
	keys_to_items map[string]*LFUCacheItem

	// number of hits
	hits int

	// number of misses
	misses int
}

// NewLFU creates a new, empty LFU.
func NewLFU(max_capacity int) *LFU {

	return &LFU{
		max_capacity:  max_capacity,
		size:          0,
		num_bindings:  0,
		keys_to_items: make(map[string]*LFUCacheItem, max_capacity),
		hits:          0,
		misses:        0,
	}
}

// Given a key, Get returns the corresponding value's size and a success boolean.
func (cache *LFU) Get(key string) (value_size int, ok bool) {

	// retrieve item associated with key
	item, ok := cache.keys_to_items[key]

	value_size = item.value_size

	if ok {
		cache.hits += 1

		// update access count of item
		item.access_count += 1

	} else {
		cache.misses += 1

		return 0, false
	}

	return value_size, ok
}

// Given a key and its value's size, Set adds them to the cache.
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
	existing_item, ok := cache.keys_to_items[key]

	if ok {
		// replace value size, update accesses, and update size of cache
		cache.size -= existing_item.value_size
		cache.size += value_length
		cache.keys_to_items[key].value_size = value_length
		cache.keys_to_items[key].access_count += 1

		// in the case that a value change for a key resulted in
		// a size overload, remove until the size meets the capacity
		for cache.size > cache.max_capacity {

			// find what key, value pair we should evict
			key_to_remove := cache.Evict_Which()

			// remove the chosen key (capacity is being updated in Remove())
			success := cache.Remove(key_to_remove)

			if success {
				cache.num_bindings -= 1
			}
		}

		return true
	}

	// if not enough space and an item with the key does not exist,
	// evict until there is enough space
	for insert_size > cache.RemainingStorage() {

		// find what key, value pair we should evict
		key_to_remove := cache.Evict_Which()

		// remove the chosen key (capacity is being updated in Remove())
		success := cache.Remove(key_to_remove)
		if success {
			cache.num_bindings -= 1
		}

	}

	// add new item with key and value length
	cache.keys_to_items[key] = &LFUCacheItem{
		value_size:   value_length,
		access_count: 1}

	// update cache fields
	cache.size += insert_size
	cache.num_bindings += 1

	return true
}

// Evict_Which() is an algorithm to select which item in the cache to evict.
func (cache *LFU) Evict_Which() (key string) {

	// make sure there are items in the cache
	if len(cache.keys_to_items) < 1 {
		return ""
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

// MaxStorage returns the maximum number of bytes this cache can store.
func (cache *LFU) MaxStorage() int {
	return cache.max_capacity
}

// RemainingStorage returns the number of unused bytes available in this cache.
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
// ok is true if a value size was found and false otherwise.
func (cache *LFU) Remove(key string) (ok bool) {

	// check if there is an item associated with key
	_, ok = cache.keys_to_items[key]
	if !ok {
		return false
	}

	// update the capacity
	cache.size -= (cache.keys_to_items[key].value_size + len([]byte(key)))

	// remove the key from the cache
	delete(cache.keys_to_items, key)

	return true
}
