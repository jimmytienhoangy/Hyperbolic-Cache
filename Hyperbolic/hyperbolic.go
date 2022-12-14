package cache

import (
	"math/rand"
	"time"
)

// A HyperbolicCacheItem is an item with metadata that
// goes in a HyperBolicCache.
type HyperbolicCacheItem struct {

	// size of item's value (bytes)
	value_size int

	// how many times the item is accessed
	access_count int

	// the date and time the item was first inserted (with its current value?)
	insert_time time.Time
}

// A HyperbolicCache is a cache that uses hyperbolic caching.
type HyperbolicCache struct {

	// maximum number of bytes the cache can hold
	max_capacity int

	// total number of bytes currently in the cache
	size int

	// number of bindings in the cache
	num_bindings int

	// map of keys to items in the cache
	keys_to_items map[string]*HyperbolicCacheItem

	// number of hits
	hits int

	// number of misses
	misses int
}

// NewHyperbolicCache creates a new, empty HyperbolicCache.
func NewHyperbolicCache(max_capacity int) *HyperbolicCache {

	return &HyperbolicCache{
		max_capacity:  max_capacity,
		size:          0,
		num_bindings:  0,
		keys_to_items: make(map[string]*HyperbolicCacheItem, max_capacity),
		hits:          0,
		misses:        0,
	}
}

// Given a key, Get returns the corresponding value's size and a success boolean.
func (cache *HyperbolicCache) Get(key string) (value_size int, ok bool) {

	// retrieve item associated with key
	item, ok := cache.keys_to_items[key]

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

// Given a key and its value's size, Set adds them to the cache.
func (cache *HyperbolicCache) Set(key string, value_size int) bool {

	// size of value to be added
	value_length := value_size

	// size of key and value size to be added
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
			_, success := cache.Remove(key_to_remove)

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
		_, success := cache.Remove(key_to_remove)
		if success {
			cache.num_bindings -= 1
		}

	}

	// add new item with key and value length
	cache.keys_to_items[key] = &HyperbolicCacheItem{
		value_size:   value_length,
		access_count: 1,
		insert_time:  time.Now()}

	// update cache fields
	cache.size += insert_size
	cache.num_bindings += 1

	return true
}

// Calc_P calculates the priority of an item for the eviction algorithm.
func (item *HyperbolicCacheItem) Calc_P() (index float32) {

	// calculate the time since item's (current value's?)
	// initial insertion into the cache
	time_in_cache := time.Now().Sub(item.insert_time)

	// priority = number of accesses / time in cache
	return float32(item.access_count) / float32(time_in_cache.Milliseconds())

}

// Evict_Which() is an algorithm to select which item in the cache to evict.
func (cache *HyperbolicCache) Evict_Which() (key string) {

	// make sure there are items in the cache
	if len(cache.keys_to_items) < 1 {
		return ""
	}

	// sample size S
	sample_size := len(cache.keys_to_items)

	// create a randomly ordered slice of the cache's current keys
	keys := make([]string, len(cache.keys_to_items))
	i := 0
	for j := range cache.keys_to_items {
		keys[i] = j
		i++
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	// take a sample of the randomly ordered slice
	sampled_items := keys[0:sample_size]

	// find the key of the sample item with the minimum p value
	minimum := sampled_items[0]
	minValue := cache.keys_to_items[sampled_items[0]].Calc_P()

	for _, key := range sampled_items {
		if cache.keys_to_items[key].Calc_P() < minValue {
			minValue = cache.keys_to_items[key].Calc_P()
			minimum = key
		}
	}

	return minimum
}

// MaxStorage returns the maximum number of bytes this cache can store.
func (cache *HyperbolicCache) MaxStorage() int {
	return cache.max_capacity
}

// RemainingStorage returns the number of unused bytes available in this cache.
func (cache *HyperbolicCache) RemainingStorage() int {
	return cache.max_capacity - cache.size
}

// Stats returns statistics about how many search hits and misses have occurred.
func (cache *HyperbolicCache) Stats() *Stats {
	return &Stats{Hits: cache.hits, Misses: cache.misses}
}

// Len returns the number of bindings in the cache.
func (cache *HyperbolicCache) Len() int {
	return cache.num_bindings
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value size was found and false otherwise.
func (cache *HyperbolicCache) Remove(key string) (value_size int, ok bool) {

	// check if there is an item associated with key
	item, ok := cache.keys_to_items[key]
	if !ok {
		return 0, false
	}

	// update the cache's current size
	cache.size -= (cache.keys_to_items[key].value_size + len([]byte(key)))

	// remove the key from the cache
	delete(cache.keys_to_items, key)

	return item.value_size, true
}
