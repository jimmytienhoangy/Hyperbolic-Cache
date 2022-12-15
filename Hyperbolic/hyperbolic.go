package cache

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

// A HyperbolicCacheItem is an item with metadata that implicitly
// holds a value of set size. It goes in a HyperbolicCache.
type HyperbolicCacheItem struct {

	// how many times the item has been accessed
	access_count int

	// when the item was first inserted
	initial_insert_time time.Time
}

// A HyperbolicCache is a cache that uses hyperbolic caching.
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
		log.Fatal("The sampling size for the hyperbolic caching algorithm can not be " +
			"greater than the number of items this cache can hold!")
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

// Get returns how many times the item associated with the given key
// has been accessed (not including this access) and a success boolean.
func (cache *HyperbolicCache) Get(key string) (access_count int, ok bool) {

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
func (cache *HyperbolicCache) Set(key string) (ok bool) {

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
	cache.keys_to_items[key] = &HyperbolicCacheItem{
		access_count:        1,
		initial_insert_time: time.Now()}

	// update size of cache
	cache.size += 1

	return true
}

// Calc_P calculates the priority of an item for the eviction algorithm.
func (item *HyperbolicCacheItem) Calc_P() (index float32) {

	// calculate the time since item's
	// initial insertion into the cache
	time_in_cache := time.Since(item.initial_insert_time)

	// priority = number of accesses / time in cache
	return float32(item.access_count) / float32(time_in_cache.Milliseconds())

}

// Evict_Which() is an algorithm to select which item in the cache to evict.
func (cache *HyperbolicCache) Evict_Which() (key string) {

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
	sampled_items := keys[0:cache.sample_size]

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

// MaxStorage returns the maximum number of items this cache can store.
func (cache *HyperbolicCache) MaxStorage() int {
	return cache.max_capacity
}

// RemainingStorage returns the number of items that can still be stored
// in this cache.
func (cache *HyperbolicCache) RemainingStorage() int {
	return cache.max_capacity - cache.size
}

// Stats returns statistics about how many search hits and misses have occurred.
func (cache *HyperbolicCache) Stats() *Stats {
	return &Stats{Hits: cache.hits, Misses: cache.misses}
}

// Len returns the number of items in the cache.
func (cache *HyperbolicCache) Len() int {
	return cache.size
}

// Remove removes the item associated with the given key from the cache, if it exists.
// ok is true if an item was found and false otherwise.
func (cache *HyperbolicCache) Remove(key string) (ok bool) {

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
