package cache

import (
	"container/list"
)

// A LRUCache is a fixed-size, in-memory cache with 
// least-recently-used eviction.
type LRUCache struct {

	// total number of items the LRUCache can store
	max_capacity int

	// total number of items currently in the LRUCache
	size int

	// mapping of keys to items in the LRUCache
	keys_to_items map[string]*list.Element

	// linked list of string keys in the LRUCache
	linked_list *list.List

	// number of hits from the LRUCache
	hits int

	// number of misses from the LRUCache
	misses int
}

// A LRUCacheItem holds a key, value pair to be put in a linked list.
type LRUCacheItem struct {
	key   string
	value int
}

// NewLRU returns a pointer to a new, empty LRUCache.
func NewLRUCache(max_capacity int) *LRUCache {

	// create and initialize a new LRUCache
	return &LRUCache{
		max_capacity:  max_capacity,
		size:          0,
		keys_to_items: make(map[string]*list.Element, max_capacity),
		linked_list:   list.New(),
		hits:          0,
		misses:        0,
	}
}

// Get returns a success boolean indicating if an item with the key was found.
// This operation counts as a "use" for that item.
func (lru *LRUCache) Get(key string) (ok bool) {

	// cache is empty
	if lru.size == 0 {
		lru.misses++
		return false
	}

	// check if there is an item with the given key
	existing_item, ok := lru.keys_to_items[key]

	// update hits/misses and possibly update most recently used
	if ok {
		lru.hits++

		// move to the back
		lru.linked_list.MoveToBack(existing_item)

	} else {
		lru.misses++
		return false
	}

	return true
}

// Remove removes an item with the given key from the LRUCache, if it exists,
// and returns a success boolean. ok is true if an item was found and false otherwise.
func (lru *LRUCache) remove(key string) (ok bool) {

	// cache is empty
	if lru.size == 0 {
		return false
	}

	// check if there is an item associated with key
	item, ok := lru.keys_to_items[key]

	if !ok {
		return false
	}

	// remove item from the mapping
	delete(lru.keys_to_items, key)

	// remove item from the linked list
	lru.linked_list.Remove(item)

	lru.size -= 1

	return true
}

// Set sets the value of the item with the given key to be the given timestamp,
// possibly evicting an item to make room for a new key insertion.
// This operation counts as a "use" for that item.
// Returns true if the item was added/updated successfully, else false.
// (Note: The assignment is actually not necessary at all.)
func (lru *LRUCache) Set(operation_timestamp int, key string) (ok bool) {

	// check if there is an existing item with the key
	existing_item, ok := lru.keys_to_items[key]

	if ok {
		existing_item.Value.(*list.Element).Value.(*LRUCacheItem).value = operation_timestamp

		// move the item to the back
		lru.linked_list.MoveToBack(existing_item)

		return true
	}

	// item with the key does not exist, so check if we need to evict
	if lru.size == lru.max_capacity {

		// remove the first item from the linked list
		first := lru.linked_list.Front()
		value := lru.linked_list.Remove(first)

		key_to_remove := value.(*list.Element).Value.(*LRUCacheItem).key

		// remove the first item from the map
		delete(lru.keys_to_items, key_to_remove)

		// update the current size of the cache
		lru.size -= 1

	}

	// insert the item into the linked list
	item_to_add := &list.Element{Value: &LRUCacheItem{key: key, value: operation_timestamp}}
	entry := lru.linked_list.PushBack(item_to_add)

	// store the mapping of the pointer of the item into the mapping
	lru.keys_to_items[key] = entry

	// update the current size of the LRUCache
	lru.size += 1

	return true
}

// Stats returns statistics about how many search hits and misses have
// occurred in the LRUCache.
func (lru *LRUCache) Stats() *Stats {
	return &Stats{Hits: lru.hits, Misses: lru.misses}
}
