package cache

// import list from Go's Standard library
import (
	"container/list"
)

// An LRU is a fixed-size, in-memory cache with least-recently-used eviction
type LRU struct {

	// total number of items the LRU can store
	max_capacity int

	// total number of items in the LRU
	current_capacity int

	// mapping of string keys to list Element pointers for constant time operations
	mapping map[string]*list.Element

	// linked list of string keys in the LRU
	linked_list *list.List

	// current number of bindings in the LRU
	num_bindings int

	// number of hits from the LRU
	hits int

	// number of misses from the LRU
	misses int
}

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {

	// create and initialize a new LRU
	return &LRU{
		max_capacity:     limit,
		current_capacity: 0,
		mapping:          make(map[string][]byte, limit),
		linked_list:      list.New(),
		num_bindings:     0,
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
	return lru.max_capacity - lru.current_capacity
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value []byte, ok bool) {

	// get the list element associated with key
	element, ok := lru.mapping[key]

	// update hits/misses and possibly update most recently used
	if ok {
		lru.hits++

		// remove from the linkedlist and add to the back
		for element := lru.linked_list.Front(); element != nil; element = element.Next() {
			if element.Value.(string) == key {
				lru.linked_list.Remove(element)
			}
		}
		lru.linked_list.PushBack(key)

	} else {
		lru.misses++
	}

	return element.Value.(*list.Element).Value.(*binding).value, ok
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value int, ok bool) {

	// get the list element associated with key
	element, ok := lru.mapping[key]

	// value associated with key does not exist
	if !ok {
		return 0, false
	}

	// remove from the hashmap
	delete(lru.mapping, key)

	// remove from the linkedlist
	for element := lru.linked_list.Front(); element != nil; element = element.Next() {
		if element.Value.(string) == key {
			lru.linked_list.Remove(element)
		}
	}

	// update current capacity and number of bindings
	// lru.current_capacity -= (len([]byte(key)) +
	// 	len(element.Value.(*list.Element).Value.(*binding).value))
	lru.current_capacity -= 1

	lru.num_bindings--

	return element.Value.(*list.Element).Value.(*binding).value, true
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(key string, value int) bool {

	// capacity to be added
	insert_capacity := 1 //(len([]byte(key)) + len(value))

	// item value pair is too large
	if (insert_capacity) > lru.max_capacity {
		return false
	}

	// check if the key already has an element
	element, ok := lru.mapping[key]
	if ok {

		// replace value, update capacity, and return
		lru.current_capacity -= len(element.Value.(*list.Element).Value.(*binding).value)
		lru.current_capacity += len(value)
		element.Value.(*list.Element).Value.(*binding).value = value

		// move to the back
		lru.linked_list.MoveToBack(element)

		// if not enough space, remove until there is enough space
		for lru.current_capacity > lru.max_capacity {

			// remove the first element
			first_test := lru.linked_list.Front()
			value := lru.linked_list.Remove(first_test)

			key_to_remove_test := value.(*list.Element).Value.(*binding).key

			// update the capacity
			// lru.current_capacity -= len(value.(*list.Element).Value.(*binding).value)
			// lru.current_capacity -= 1

			// remove the first value from the map
			delete(lru.mapping, key_to_remove_test)
			// lru.num_bindings--
		}
		return true
	}

	// if not enough space, remove until there is enough space
	for insert_capacity > lru.RemainingStorage() {

		// remove the first value
		first := lru.linked_list.Front()
		value := lru.linked_list.Remove(first)

		key_to_remove := value.(*list.Element).Value.(*binding).key

		// update the capacity
		// lru.current_capacity -= len(value.(*list.Element).Value.(*binding).value)
		lru.current_capacity -= 1 //len((key_to_remove))

		// remove the first value from the map
		delete(lru.mapping, key_to_remove)
		lru.num_bindings--
	}

	// insert the value into the linkedlist
	element_to_add := &list.Element{Value: &binding{key: key, value: value}}
	entry := lru.linked_list.PushBack(element_to_add)

	// store the mapping of the pointer of the element into the hashmap
	lru.mapping[key] = entry

	// update current capacity and binding count
	// lru.current_capacity += insert_capacity
	lru.current_capacity += 1
	lru.num_bindings++

	return true
}

// Len returns the number of bindings in the LRU
func (lru *LRU) Len() int {
	return lru.num_bindings
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lru *LRU) Stats() *Stats {
	return &Stats{Hits: lru.hits, Misses: lru.misses}
}
