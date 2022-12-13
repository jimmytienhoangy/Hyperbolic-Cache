package cache

// import list from Go's Standard library
import (
	"container/list"
)

// A FIFO is a fixed-size, in-memory cache with first-in first-out eviction
type HyperbolicCache struct {

	// total number of bytes the FIFO can store
	max_capacity int

	// total number of current bytes in the FIFO
	current_capacity int

	// mapping of keys to values for constant time operations
	mapping map[string][]byte

	// linked list of string keys in the FIFO
	linked_list *list.List

	// current number of bindings in the FIFO
	num_bindings int

	// number of hits from the FIFO
	hits int

	// number of misses from the FIFO
	misses int
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewHyperbolicCache(limit int) *HyperbolicCache {

	// create and initialize a new FIFO
	return &HyperbolicCache{
		max_capacity:     limit,
		current_capacity: 0,
		mapping:          make(map[string][]byte, limit),
		linked_list:      list.New(),
		num_bindings:     0,
		hits:             0,
		misses:           0,
	}
}

// MaxStorage returns the maximum number of bytes this FIFO can store
func (cache *HyperbolicCache) MaxStorage() int {
	return cache.max_capacity
}

// RemainingStorage returns the number of unused bytes available in this FIFO
func (cache *HyperbolicCache) RemainingStorage() int {
	return cache.max_capacity - cache.current_capacity
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (cache *HyperbolicCache) Get(key string) (value []byte, ok bool) {

	// get the value associated with key
	value, ok = cache.mapping[key]

	// update hits/misses
	if ok {
		cache.hits++
	} else {
		cache.misses++
	}

	return value, ok
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (cache *HyperbolicCache) Remove(key string) (value []byte, ok bool) {

	// get the value associated with key
	value, ok = cache.mapping[key]

	// value associated with key does not exist
	if !ok {
		return nil, false
	}

	// remove from the hashmap
	delete(cache.mapping, key)

	// remove from the linkedlist
	for element := cache.linked_list.Front(); element != nil; element = element.Next() {
		if element.Value.(string) == key {
			cache.linked_list.Remove(element)
		}
	}

	// update current capacity and number of bindings
	cache.current_capacity -= (len([]byte(key)) + len(value))
	cache.num_bindings--

	return value, true
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (cache *HyperbolicCache) Set(key string, value []byte) bool {

	// capacity to be added
	insert_capacity := (len([]byte(key)) + len(value))

	// item value pair is too large
	if (insert_capacity) > cache.max_capacity {
		return false
	}

	// check if the key already has a value
	existing_value, ok := cache.mapping[key]
	if ok {
		// replace value, update capacity, and return
		cache.current_capacity -= len(existing_value)
		cache.current_capacity += len(value)
		cache.mapping[key] = value
		return true
	}

	// if not enough space, remove until there is enough space
	for insert_capacity > cache.RemainingStorage() {

		// remove the first value
		first := cache.linked_list.Front()
		cache.linked_list.Remove(first)

		key_to_remove := first.Value.(string)

		// update the capacity
		cache.current_capacity -= len(cache.mapping[key_to_remove])
		cache.current_capacity -= len([]byte(key_to_remove))

		// remove the first value from the map
		delete(cache.mapping, key_to_remove)
		cache.num_bindings--
	}

	// insert the key into the linkedlist and update the corresponding hashmap value
	cache.linked_list.PushBack(key)
	cache.mapping[key] = value

	// update the capacity and nubmer of bindings
	cache.current_capacity += insert_capacity
	cache.num_bindings++

	return true
}

// Len returns the number of bindings in the FIFO.
func (cache *HyperbolicCache) Len() int {
	return cache.num_bindings
}

// Stats returns statistics about how many search hits and misses have occurred.
func (cache *HyperbolicCache) Stats() *Stats {
	return &Stats{Hits: cache.hits, Misses: cache.misses}
}

// Don't know how I feel about stats here
