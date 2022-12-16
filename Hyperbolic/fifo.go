package cache

import (
	"container/list"
)

// A FIFOCache is a fixed-size, in-memory cache with 
// first-in, first-out eviction.
type FIFOCache struct {

	// total number of items the FIFOCache can store
	max_capacity int

	// total number of items currently in the FIFOCache
	size int

	// mapping of keys to items (represented by an arbitrary int)
	// in the FIFOCache
	keys_to_items map[string]int

	// linked list of string keys in the FIFOCache
	linked_list *list.List

	// number of hits from the FIFOCache
	hits int

	// number of misses from the FIFOCache
	misses int
}

// NewFIFOCache returns a pointer to a new, empty FIFOCache.
func NewFIFOCache(max_capacity int) *FIFOCache {

	// create and initialize a new FIFOCache
	return &FIFOCache{
		max_capacity:  max_capacity,
		size:          0,
		keys_to_items: make(map[string]int, max_capacity),
		linked_list:   list.New(),
		hits:          0,
		misses:        0,
	}
}

// Get returns a success boolean indicating if an item with the
// key was found in the cache.
func (fifo *FIFOCache) Get(key string) (ok bool) {

	// cache is empty
	if fifo.size == 0 {
		fifo.misses++
		return false
	}

	// check if there is an item with the given key
	_, ok = fifo.keys_to_items[key]

	// update hits/misses
	if ok {
		fifo.hits++
	} else {
		fifo.misses++
		return false
	}

	return true
}

// Remove removes the item associated with the given key from
// the FIFOCache, if it exists. ok is true if an item was found
// and false otherwise.
func (fifo *FIFOCache) remove(key string) (ok bool) {

	// cache is empty
	if fifo.size == 0 {
		return false
	}

	// check if there is an item associated with key
	_, ok = fifo.keys_to_items[key]

	if !ok {
		return false
	}

	// remove the item from the mapping
	delete(fifo.keys_to_items, key)

	// remove the item from the linked list
	for element := fifo.linked_list.Front(); element != nil; element = element.Next() {
		if element.Value.(string) == key {
			fifo.linked_list.Remove(element)
		}
	}

	// update the current size of the FIFOCache
	fifo.size -= 1

	return true
}

// Set sets the value of the item with the given key to be the given timestamp,
// possibly evicting an item to make room for a new key insertion.
// Returns true if the item was added/updated successfully, else false.
// (Note: The assignment is actually not necessary at all.)
func (fifo *FIFOCache) Set(operation_timestamp int, key string) (ok bool) {

	// check if there is an existing item with the key
	_, ok = fifo.keys_to_items[key]

	if ok {
		fifo.keys_to_items[key] = operation_timestamp
		return true
	}

	// item with the key does not exist, so check if we need to evict
	if fifo.size == fifo.max_capacity {

		// remove the first item
		first := fifo.linked_list.Front()
		fifo.linked_list.Remove(first)

		key_to_remove := first.Value.(string)

		// remove the first item from the map
		delete(fifo.keys_to_items, key_to_remove)

		fifo.size--
	}

	// insert the value into the linked list
	fifo.linked_list.PushBack(key)
	fifo.keys_to_items[key] = operation_timestamp

	// update the size of the FIFOCache
	fifo.size++

	return true
}

// Stats returns statistics about how many search hits and misses have occurred.
func (fifo *FIFOCache) Stats() *Stats {
	return &Stats{Hits: fifo.hits, Misses: fifo.misses}
}
