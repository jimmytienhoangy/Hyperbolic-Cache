package cache

// import list from Go's Standard library
import (
	"container/list"
)

// A FIFO is a fixed-size, in-memory cache with first-in first-out eviction
type FIFO struct {

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
func NewFifo(limit int) *FIFO {

	// create and initialize a new FIFO
	return &FIFO{
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
func (fifo *FIFO) MaxStorage() int {
	return fifo.max_capacity
}

// RemainingStorage returns the number of unused bytes available in this FIFO
func (fifo *FIFO) RemainingStorage() int {
	return fifo.max_capacity - fifo.current_capacity
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (fifo *FIFO) Get(key string) (value []byte, ok bool) {

	// get the value associated with key
	value, ok = fifo.mapping[key]

	// update hits/misses
	if ok {
		fifo.hits++
	} else {
		fifo.misses++
	}

	return value, ok
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {

	// get the value associated with key
	value, ok = fifo.mapping[key]

	// value associated with key does not exist
	if !ok {
		return nil, false
	}

	// remove from the hashmap
	delete(fifo.mapping, key)

	// remove from the linkedlist
	for element := fifo.linked_list.Front(); element != nil; element = element.Next() {
		if element.Value.(string) == key {
			fifo.linked_list.Remove(element)
		}
	}

	// update current capacity and number of bindings
	fifo.current_capacity -= (len([]byte(key)) + len(value))
	fifo.num_bindings--

	return value, true
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(key string, value []byte) bool {

	// capacity to be added
	insert_capacity := (len([]byte(key)) + len(value))

	// item value pair is too large
	if (insert_capacity) > fifo.max_capacity {
		return false
	}

	// check if the key already has a value
	existing_value, ok := fifo.mapping[key]
	if ok {
		// replace value, update capacity, and return
		fifo.current_capacity -= len(existing_value)
		fifo.current_capacity += len(value)
		fifo.mapping[key] = value
		return true
	}

	// if not enough space, remove until there is enough space
	for insert_capacity > fifo.RemainingStorage() {

		// remove the first value
		first := fifo.linked_list.Front()
		fifo.linked_list.Remove(first)

		key_to_remove := first.Value.(string)

		// update the capacity
		fifo.current_capacity -= (len(fifo.mapping[key_to_remove]) + len([]byte(key_to_remove)))

		// remove the first value from the map
		delete(fifo.mapping, key_to_remove)
		fifo.num_bindings--
	}

	// insert the value into the linkedlist
	fifo.linked_list.PushBack(key)
	fifo.mapping[key] = value

	// update the capacity and nubmer of bindings
	fifo.current_capacity += insert_capacity
	fifo.num_bindings++

	return true
}

// Len returns the number of bindings in the FIFO.
func (fifo *FIFO) Len() int {
	return fifo.num_bindings
}

// Stats returns statistics about how many search hits and misses have occurred.
func (fifo *FIFO) Stats() *Stats {
	return &Stats{Hits: fifo.hits, Misses: fifo.misses}
}
