package cache

import (
	"container/list"
)

// A LRUCache is a fixed-size, in-memory cache with least-recently-used eviction.
type LRUCache struct {

	// total number of items the LRUCache can store
	max_capacity int

	// total number of items currently in the LRUCache
	size int

	// mapping of keys to items in the LRUCache
	mapping map[string]*list.Element

	// linked list of string keys in the LRUCache
	linked_list *list.List

	// number of hits from the LRUCache
	hits int

	// number of misses from the LRUCache
	misses int
}

// A binding holds a key value pair to be put in a linkedlist
type binding struct {
	key   string
	value int
}

// NewLRU returns a pointer to a new, empty LRUCache.
func NewLRU(max_capacity int) *LRUCache {

	// create and initialize a new LRU
	return &LRUCache{
		max_capacity:     max_capacity,
		size: 0,
		mapping:          make(map[string]*list.Element, max_capacity),
		linked_list:      list.New(),
		hits:             0,
		misses:           0,
	}
}

// MaxStorage returns the maximum number of bytes this LRU can store
func (lru *LRU) MaxStorage() int {
	return lru.max_capacity
}

// RemainingStorage returns the number of unused bytes available in this LRU
func (lru *LRU) RemainingStorage() int {
	return lru.max_capacity - lru.size
}

// Get "returns" the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value int, ok bool) {

	// get the list element associated with key
	element, ok := lru.mapping[key]

	// update hits/misses and possibly update most recently used
	if ok {
		lru.hits++

		// move to the back
		lru.linked_list.MoveToBack(element)

	} else {
		lru.misses++
		//return 0, ok
	}

	//return element.Value.(*list.Element).Value.(*binding).value, ok
	return 0, ok
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (ok bool) {

	// get the list element associated with key
	element, ok := lru.mapping[key]

	// value associated with key does not exist
	if !ok {
		return false
	}

	// remove from the hashmap
	delete(lru.mapping, key)

	// remove from the linkedlist
	lru.linked_list.Remove(element)

	// update current capacity and number of bindings
	// lru.size -= (len([]byte(key)) +
	// 	len(element.Value.(*list.Element).Value.(*binding).value))
	lru.size -= 1

	//return element.Value.(*list.Element).Value.(*binding).value, true
	return true
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(timestamp int, key string) bool {

	// capacity to be added
	// insert_capacity := (len([]byte(key)) + len(value))
	//insert_capacity := 1

	// item value pair is too large
	//if (insert_capacity) > lru.max_capacity {
	//	return false
	//}

	// check if the key already has an element
	element, ok := lru.mapping[key]
	if ok {

		// replace value, update capacity, and return
		// lru.size -= len(element.Value.(*list.Element).Value.(*binding).value)
		// lru.size += len(value)
		element.Value.(*list.Element).Value.(*binding).value = 1

		// move to the back
		lru.linked_list.MoveToBack(element)

		// // if not enough space, remove until there is enough space
		// for lru.size > lru.max_capacity {

		// 	// remove the first element
		// 	first_test := lru.linked_list.Front()
		// 	value := lru.linked_list.Remove(first_test)

		// 	key_to_remove_test := value.(*list.Element).Value.(*binding).key

		// 	// update the capacity
		// 	lru.size -= len(value.(*list.Element).Value.(*binding).value)
		// 	lru.size -= len((key_to_remove_test))

		// 	// remove the first value from the map
		// 	delete(lru.mapping, key_to_remove_test)
		// 	lru.num_bindings--
		// }
		return true
	}

	// if not enough space, remove until there is enough space
	//for insert_capacity > lru.RemainingStorage() {
	
		if lru.size == lru.max_capacity {

		// remove the first value
		first := lru.linked_list.Front()
		value := lru.linked_list.Remove(first)

		key_to_remove := value.(*list.Element).Value.(*binding).key

		// update the capacity
		lru.size -= 1

		// remove the first value from the map
		delete(lru.mapping, key_to_remove)
	}

	// insert the value into the linkedlist
	element_to_add := &list.Element{Value: &binding{key: key, value: 1}}
	entry := lru.linked_list.PushBack(element_to_add)

	// store the mapping of the pointer of the element into the hashmap
	lru.mapping[key] = entry

	// update current capacity and binding count
	lru.size += 1

	return true
}

// Len returns the number of items in the LRU.
func (lru *LRU) Len() int {
	return lru.size
}

// Stats returns statistics about how many search hits and misses have 
// occurred in the LRU.
func (lru *LRU) Stats() *Stats {
	return &Stats{Hits: lru.hits, Misses: lru.misses}
}
