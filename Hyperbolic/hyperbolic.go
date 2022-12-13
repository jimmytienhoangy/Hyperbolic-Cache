package cache

import (
	//"container/list"
	"math/rand"
	"time"
)

// items that will go in cache
type CacheItem struct {
	// identifies item
	key string
	// value of item
	value []byte
	// how many times the item is accessed
	access_count int
	// the date and time the item was first inserted
	insert_time time.Time
}

//
type HyperbolicCache struct {
	// max number of items the cache can hold
	max_capacity int

	// current number of items in the cache
	size int

	// map from keys to items in cache
	mapping map[string]*CacheItem
}

// creates a new, empty Hyperbolic cache
func NewHyperbolicCache(max_capacity int) *HyperbolicCache {

	return &HyperbolicCache{
		max_capacity: max_capacity,
		size:         0,
		mapping:      make(map[string]*CacheItem, max_capacity),
	}
}

// given a key, return the value and a boolean
// func (cache *HyperbolicCache) Get(key string) (value interface{}, ok bool) {
func (cache *HyperbolicCache) Get(key string) (value []byte, ok bool) {
	item, ok := cache.mapping[key]
	item.access_count += 1
	value = item.value

	return value, ok
}

// func (cache *HyperbolicCache) Set(key string, value interface{}) bool {
func (cache *HyperbolicCache) Set(key string, value []byte) bool {

	// check for max capacity first

	// if cache is at max capacity, evict

	// otherwise, insert into the cache

	// do the eviction
	if cache.size == cache.max_capacity {
		var item_to_evict string = cache.Evict_which() // key that will be evicted
		delete(cache.mapping, item_to_evict)
	}

	// some sort of return false here

	//otherwise, insert to cache

	// is this how we do it? 
	cache.mapping[key] = &CacheItem{key: key, value: value, access_count: 0, insert_time: time.Now()}

	return true
}

// maybe not have a pointer method?
//func (item *CacheItem) calc_p() (index float32) {
func (item *CacheItem) calc_p() (index float32) {
	time_in_cache := time.Now().Sub(item.insert_time)

	return float32(item.access_count) / float32(time_in_cache.Milliseconds()) // I converted it to float

}

// Note: This assumes all keys are unique:
func (cache *HyperbolicCache) Evict_which() (key string) {
	//sampled_items := rand.Perm(cache.size)[0:5] // generate a random order of this sizehttps://golangbyexample.com/generate-random-array-slice-golang/

	// get a randomly ordered slice of keys
	keys := make([]string, len(cache.mapping))

	i := 0
	for k := range cache.mapping {
		keys[i] = k
		i++
	}

	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	sampled_items := keys[0:5]

	// get the minimum
	minimum := keys[0]
	minvalue := cache.mapping[keys[0]].calc_p()

	for _, key := range sampled_items {
		if cache.mapping[key].calc_p() < minvalue {
			minvalue = cache.mapping[key].calc_p()
			minimum = key
		}
	}

	return minimum
}

// sampled_items = random_sample(S)
// return argmin(p(i) for i in sampled_items)

// Don't know how I feel about stats here

//Additional Code

// MaxStorage returns the maximum number of bytes this FIFO can store
func (cache *HyperbolicCache) MaxStorage() int {
	return cache.max_capacity
}

// RemainingStorage returns the number of unused bytes available in this FIFO
func (cache *HyperbolicCache) RemainingStorage() int {
	return cache.max_capacity - cache.size
}

// Stats returns statistics about how many search hits and misses have occurred.
func (cache *HyperbolicCache) Stats() *Stats {
}

// Len returns the number of bindings in the FIFO.
func (cache *HyperbolicCache) Len() int {
	return cache.size
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (cache *HyperbolicCache) Remove(key string) (value []byte, ok bool) {

}
