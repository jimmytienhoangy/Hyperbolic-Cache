package cache

import (
	"log"
	"strconv"
)

// A HyperbolicCacheItem is an item with metadata that implicitly
// holds a value of set size. It goes in a HyperbolicCache.
type HyperbolicCacheItem struct {

	// how many times the item has been accessed
	access_count int

	// when the item was first inserted
	initial_timestamp int
}

// A HyperbolicCache is a cache that uses the hyperbolic 
// caching algorithm.
type HyperbolicCache struct {

	// maximum number of items the cache can hold
	max_capacity int

	// total number of items currently in the cache
	size int

	// map of keys to items in the cache
	keys_to_items map[string]*HyperbolicCacheItem

	// sample size for eviction
	sample_size int

	// number of hits
	hits int

	// number of misses
	misses int
}

// NewHyperbolicCache creates a new, empty HyperbolicCache.
func NewHyperbolicCache(max_capacity int, sample_size int) *HyperbolicCache {

	if sample_size > max_capacity {
		log.Fatal("The sampling size for the hyperbolic caching " + 
		"algorithm can not be greater than the number of items " +
		"this cache can hold!")
	}

	return &HyperbolicCache{
		max_capacity:  max_capacity,
		size:          0,
		keys_to_items: make(map[string]*HyperbolicCacheItem, max_capacity),
		sample_size:   sample_size,
		hits:          0,
		misses:        0,
	}
}

// Get returns a success boolean indicating if an item 
// with the key was found.
func (cache *HyperbolicCache) Get(key string) (ok bool) {

	// retrieve item associated with key
	item, ok := cache.keys_to_items[key]

	if ok {
		cache.hits += 1

		// update access count of item
		item.access_count += 1

	} else {
		cache.misses += 1

		return false
	}

	return true
}

// Set adds/updates an item with the given key in the cache
// and returns a success boolean.
func (cache *HyperbolicCache) Set(operation_timestamp int, key string) (ok bool) {

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

		key_to_remove := cache.evict_Which(operation_timestamp)

		success := cache.remove(key_to_remove)
		if !success {
			log.Fatal("Failed to evict an item.")
		}

	}

	// add new item with key
	cache.keys_to_items[key] = &HyperbolicCacheItem{
		access_count:      1,
		initial_timestamp: operation_timestamp}

	// update size of cache
	cache.size += 1

	return true
}

// calc_P calculates the priority of an item for the eviction algorithm.
func (item *HyperbolicCacheItem) calc_P(eviction_timestamp int) (priority float32) {

	// calculate the time since item's
	// initial insertion into the cache
	time_in_cache := eviction_timestamp - item.initial_timestamp

	// priority = number of accesses / time in cache
	return float32(item.access_count) / float32(time_in_cache)

}

// evict_Which() is an algorithm to select which item in the cache to evict.
func (cache *HyperbolicCache) evict_Which(eviction_timestamp int) (key string) {

	// make sure cache is actually full before evicting
	if cache.size != cache.max_capacity {
		log.Fatal("Should not be evicting when cache is not full.")
	}

	// make sure there are enough items in the cache to sample
	if cache.size < cache.sample_size {
		log.Fatal("Not enough items in the cache to take a sample of size " +
			strconv.Itoa(cache.sample_size) + "!")
	}

	// create a randomly ordered slice of the cache's current keys
	// iteration over maps is random in golang
	random_sample_keys := make([]string, cache.sample_size)
	count := 0
	for random_key := range cache.keys_to_items {
		random_sample_keys[count] = random_key
		count++
		if count == cache.sample_size {
			break
		}
	}

	// find the key of the sample item with the minimum p value
	minimum := random_sample_keys[0]
	minValue := cache.keys_to_items[random_sample_keys[0]].calc_P(eviction_timestamp)

	for _, key := range random_sample_keys {
		if cache.keys_to_items[key].calc_P(eviction_timestamp) < minValue {
			minValue = cache.keys_to_items[key].calc_P(eviction_timestamp)
			minimum = key
		}
	}

	return minimum
}

// Stats returns statistics about how many search hits and misses have occurred.
func (cache *HyperbolicCache) Stats() *Stats {
	return &Stats{Hits: cache.hits, Misses: cache.misses}
}

// Remove removes the item associated with the given key from the cache, if it exists.
// ok is true if an item was found and false otherwise.
func (cache *HyperbolicCache) remove(key string) (ok bool) {

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
