package cache

import (
	"container/list"
)

// A FIFOCacheItem is an item that implicitly
// holds a value of set size. It goes in a FIFOCache.
type FIFOCacheItem struct {

}

// A FIFOCache is a fixed-size, in-memory cache with first-in, first-out eviction.
type FIFOCache struct {

	// total number of items the FIFOCache can store
	max_capacity int

	// total number of items currently in the FIFOCache
	size int

	// mapping of keys to items
	keys_to_items map[string]*FIFOCacheItem

	// linked list of string keys in the FIFO
	linked_list *list.List

	// number of hits from the FIFO
	hits int

	// number of misses from the FIFO
	misses int
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit items
func NewFifo(max_capacity int) *FIFO {

	// create and initialize a new FIFO
	return &FIFO{
		max_capacity: max_capacity,
		size:         0,
		keys_to_items:      make(map[string]*FIFOCacheItem, max_capacity),
		linked_list:  list.New(),
		hits:         0,
		misses:       0,
	}
}

// MaxStorage returns the maximum number of items this FIFO can store
func (fifo *FIFO) MaxStorage() int {
	return fifo.max_capacity
}

// RemainingStorage returns the remaining number of items that can be stored in this FIFO
func (fifo *FIFO) RemainingStorage() int {
	return fifo.max_capacity - fifo.size
}

// Get "returns" the access count of the item associated with the 
// given key, if it exists. ok is true if a value was found and false otherwise.
func (fifo *FIFO) Get(key string) (access_count int, ok bool) {
	
	// access_count doesn't actually matter

	// get the value associated with key
	_, ok = fifo.keys_to_items[key]

	// update hits/misses
	if ok {
		fifo.hits++
	} else {
		fifo.misses++
	}

	return 0, ok
}

// Remove removes the item associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (ok bool) {

	// get the value associated with key
	_, ok = fifo.keys_to_items[key]

	// value associated with key does not exist
	if !ok {
		return false
	}

	// remove from the hashmap
	delete(fifo.keys_to_items, key)

	// remove from the linkedlist
	for element := fifo.linked_list.Front(); element != nil; element = element.Next() {
		if element.Value.(string) == key {
			fifo.linked_list.Remove(element)
		}
	}

	// update current capacity and number of bindings
	fifo.size -= 1

	return true
}

// Set associates the given value with the given key, possibly evicting items
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(timestamp int, key string) bool {

	// check if the key already has a value
	_, ok := fifo.keys_to_items[key]
	if ok {
		// replace value, update capacity, and return
		return true
	}

	// if not enough space, remove until there is enough space
	for fifo.size == fifo.max_capacity {

		// remove the first value
		first := fifo.linked_list.Front()
		fifo.linked_list.Remove(first)

		key_to_remove := first.Value.(string)

		// update the capacity
		fifo.size -= 1

		// remove the first value from the map
		delete(fifo.keys_to_items, key_to_remove)
	}

	// insert the value into the linkedlist
	fifo.linked_list.PushBack(key)
	fifo.keys_to_items[key] = &FIFOCacheItem{}

	// update the capacity and nubmer of bindings
	fifo.size += 1

	return true
}

// Len returns the number of items in the FIFO.
func (fifo *FIFO) Len() int {
	return fifo.size
}

// Stats returns statistics about how many search hits and misses have 
// occurred in the FIFO.
func (fifo *FIFO) Stats() *Stats {
	return &Stats{Hits: fifo.hits, Misses: fifo.misses}
}
