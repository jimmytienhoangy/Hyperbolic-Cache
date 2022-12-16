/******************************************************************************
 * caching_unit_test.go
 * Authors: Reuben Agogoe, Stephen Dong, Jimmy Hoang
 * Usage: `go test`  or  `go test -v`
 * Description: A unit testing suite for our caches.
 ******************************************************************************/

package cache

import (
	"fmt"
	"testing"
)

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

// Tests the creation of a hyperbolic cache. Then performs set and get operations.
func Test_CreateHyperbolic(t *testing.T) {
	max_capacity := 50
	sample_size := 50
	hyperbolic := NewHyperbolicCache(max_capacity, sample_size)
	for i := 0; i < max_capacity; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := hyperbolic.Set(0, key)
		if !set_success {
			t.Errorf("Failed to set binding with key: %s", key)
			t.FailNow()
		}
		get_success := hyperbolic.Get(key)
		if !get_success {
			t.Errorf("Failed to get binding with key: %s", key)
			t.FailNow()
		}
	}
}

// Tests a hyperbolic cache with a max capacity of 0 items.
func Test_EmptyHyperbolic(t *testing.T) {
	max_capacity := 0
	hyperbolic := NewHyperbolicCache(max_capacity, 0)

	// try to set bindings
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := hyperbolic.Set(i, key)

		if set_success {
			t.Errorf("This set operation should have failed!")
			t.FailNow()
		}

		get_success := hyperbolic.Get(key)
		if get_success {
			t.Errorf("This get operation should have failed!")
			t.FailNow()
		}
	}
}

// Checks to see that eviction is occuring at all.
func Test_HyperbolicEviction(t *testing.T) {
	max_capacity := 3
	sample_size := 3
	hyperbolic := NewHyperbolicCache(max_capacity, sample_size)

	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := hyperbolic.Set(i, key)
		if !set_success {
			t.Errorf("Failed to set binding with key: %s", key)
			t.FailNow()
		}
	}

	set_success := hyperbolic.Set(5, "A")
	if !set_success {
		t.Errorf("Failed to set binding with key: %s", "A")
		t.FailNow()
	}

	get_success := hyperbolic.Get("0")
	if get_success {
		t.Errorf("Item with key '0' should have been evicted.")
		t.FailNow()
	}
}

// Tests the accuracy of the hyperbolic function.
func CheckHyperbolicFunction(t *testing.T) {

	max_capacity := 5

	// sample size is the same as max capacity to make sure we
	// calculate the priority of all items (so that we can
	// check accuracy)
	hyperbolic := NewHyperbolicCache(max_capacity, max_capacity)

	var test_values [5]string
	test_values[0] = "a"
	test_values[1] = "b"
	test_values[2] = "c"
	test_values[3] = "d"
	test_values[4] = "e"

	// set bindings
	for i := 0; i < 5; i++ {
		key := test_values[i]
		set_success := hyperbolic.Set(i, key)
		if !set_success {
			t.Errorf("Failed to set binding with key: %s", key)
			t.FailNow()
		}
	}

	// get each set item
	hyperbolic.Get("e")
	hyperbolic.Get("d")
	hyperbolic.Get("c")
	hyperbolic.Get("b")
	hyperbolic.Get("a")

	// try to set when the cache is full
	hyperbolic.Set(5, "f")

	// item with key 'a' should be evicted
	for keys := range hyperbolic.keys_to_items {
		fmt.Println(keys)
	}

	get_success1 := hyperbolic.Get("f")

	if !get_success1 {
		t.Errorf("Item with key 'f' should be in the cache!")
		t.FailNow()
	}

	get_success2 := hyperbolic.Get("a")

	if get_success2 {
		t.Errorf("Item with key 'a' should have been evicted!")
		t.FailNow()
	}
}

// Test whether the hyperbolic cache is working (no setting).
func TestHyperbolicFunction2(t *testing.T) {
	capacity := 5
	cache := NewHyperbolicCache(capacity, capacity)

	var values [5]string
	values[0] = "a"
	values[1] = "b"
	values[2] = "c"
	values[3] = "d"
	values[4] = "e"

	for i := 0; i < 5; i++ {
		key := values[i]
		timestamp := (i + 1) * 2
		ok := cache.Set(timestamp, key)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}
	}

	cache.Set((5+1)*2, "f")

	ok := cache.Get("a")

	if ok {
		t.Errorf("a should have been evicted.")
		t.FailNow()
	}
}

// Test whether the hyperbolic cache is working (evicting newly entered).
func TestHyperbolicFunction3(t *testing.T) {
	capacity := 5
	cache := NewHyperbolicCache(capacity, capacity)

	var values [5]string
	values[0] = "a"
	values[1] = "b"
	values[3] = "d"
	values[4] = "e"

	// set bindings
	for i := 0; i < 5; i++ {
		key := values[i]
		timestamp := (i + 1) * 2
		ok := cache.Set(timestamp, key)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}
	}

	for i := 0; i < 5; i++ {
		cache.Get("a")
		cache.Get("b")
		cache.Get("c")
		cache.Get("d")
	}

	cache.Set((5+1)*2, "f")

	ok := cache.Get("e")

	if ok {
		t.Errorf("e should have been evicted.")
		t.FailNow()
	}
}

/*********************************************************************/

// Tests the creation of a FIFO cache. Then performs set and get operations.
func Test_CreateFIFO(t *testing.T) {
	max_capacity := 50
	fifo := NewFIFOCache(max_capacity)
	for i := 0; i < max_capacity; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := fifo.Set(0, key)
		if !set_success {
			t.Errorf("Failed to set binding with key: %s", key)
			t.FailNow()
		}
		get_success := fifo.Get(key)
		if !get_success {
			t.Errorf("Failed to get binding with key: %s", key)
			t.FailNow()
		}
	}
}

// Tests a FIFO cache with a max capacity of 0 items.
func Test_EmptyFIFO(t *testing.T) {
	max_capacity := 0
	fifo := NewFIFOCache(max_capacity)

	// try to set bindings
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := fifo.Set(0, key)

		if set_success {
			t.Errorf("This set operation should have failed!")
			t.FailNow()
		}

		get_success := fifo.Get(key)
		if get_success {
			t.Errorf("This get operation should have failed!")
			t.FailNow()
		}
	}
}

// Checks to see that FIFO eviction is occuring at all and accurately.
func Test_FIFOEviction(t *testing.T) {
	max_capacity := 3
	fifo := NewFIFOCache(max_capacity)

	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := fifo.Set(0, key)
		if !set_success {
			t.Errorf("Failed to set binding with key: %s", key)
			t.FailNow()
		}
	}

	set_success := fifo.Set(0, "A")
	if !set_success {
		t.Errorf("Failed to set binding with key: %s", "A")
		t.FailNow()
	}

	get_success := fifo.Get("0")
	if get_success {
		t.Errorf("Item with key '0' should have been evicted.")
		t.FailNow()
	}

	set_success2 := fifo.Set(0, "B")
	if !set_success2 {
		t.Errorf("Failed to set binding with key: %s", "B")
		t.FailNow()
	}

	get_success2 := fifo.Get("1")
	if get_success2 {
		t.Errorf("Item with key '1' should have been evicted.")
		t.FailNow()
	}
}

/*********************************************************************/

// Tests the creation of a LRU cache. Then performs set and get operations.
func Test_CreateLRU(t *testing.T) {
	max_capacity := 50
	lru := NewLRUCache(max_capacity)
	for i := 0; i < max_capacity; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := lru.Set(0, key)
		if !set_success {
			t.Errorf("Failed to set binding with key: %s", key)
			t.FailNow()
		}
		get_success := lru.Get(key)
		if !get_success {
			t.Errorf("Failed to get binding with key: %s", key)
			t.FailNow()
		}
	}
}

// Tests a LRU cache with a max capacity of 0 items.
func Test_EmptyLRU(t *testing.T) {
	max_capacity := 0
	lru := NewLRUCache(max_capacity)

	// try to set bindings
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := lru.Set(0, key)

		if set_success {
			t.Errorf("This set operation should have failed!")
			t.FailNow()
		}

		get_success := lru.Get(key)
		if get_success {
			t.Errorf("This get operation should have failed!")
			t.FailNow()
		}
	}
}

// Checks to see that LRU eviction is occuring at all and accurately.
func Test_LRUEviction(t *testing.T) {
	max_capacity := 3
	lru := NewLRUCache(max_capacity)

	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := lru.Set(0, key)
		if !set_success {
			t.Errorf("Failed to set binding with key: %s", key)
			t.FailNow()
		}
	}

	lru.Get("0")

	set_success := lru.Set(0, "A")
	if !set_success {
		t.Errorf("Failed to set binding with key: %s", "A")
		t.FailNow()
	}

	get_success := lru.Get("1")
	if get_success {
		t.Errorf("Item with key '1' should have been evicted.")
		t.FailNow()
	}

	set_success2 := lru.Set(0, "B")
	if !set_success2 {
		t.Errorf("Failed to set binding with key: %s", "B")
		t.FailNow()
	}

	get_success2 := lru.Get("2")
	if get_success2 {
		t.Errorf("Item with key '2' should have been evicted.")
		t.FailNow()
	}
}

/*********************************************************************/

// Tests the creation of a LFU cache. Then performs set and get operations.
func Test_CreateLFU(t *testing.T) {
	max_capacity := 50
	lfu := NewLRUCache(max_capacity)
	for i := 0; i < max_capacity; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := lfu.Set(0, key)
		if !set_success {
			t.Errorf("Failed to set binding with key: %s", key)
			t.FailNow()
		}
		get_success := lfu.Get(key)
		if !get_success {
			t.Errorf("Failed to get binding with key: %s", key)
			t.FailNow()
		}
	}
}

// Tests a LFU cache with a max capacity of 0 items.
func Test_EmptyLFU(t *testing.T) {
	max_capacity := 0
	lfu := NewLRUCache(max_capacity)

	// try to set bindings
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := lfu.Set(0, key)

		if set_success {
			t.Errorf("This set operation should have failed!")
			t.FailNow()
		}

		get_success := lfu.Get(key)
		if get_success {
			t.Errorf("This get operation should have failed!")
			t.FailNow()
		}
	}
}

// Checks to see that LFU eviction is occuring at all and accurately.
func Test_LFUEviction(t *testing.T) {
	max_capacity := 5
	lfu := NewLFUCache(max_capacity)

	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("%d", i)
		set_success := lfu.Set(0, key)
		if !set_success {
			t.Errorf("Failed to set binding with key: %s", key)
			t.FailNow()
		}
	}

	lfu.Get("0")
	lfu.Get("0")
	lfu.Get("0")
	lfu.Get("0")
	lfu.Get("0")
	lfu.Set(0, "1")
	lfu.Get("2")
	lfu.Get("1")
	lfu.Get("3")
	lfu.Set(0, "3")
	lfu.Get("3")
	lfu.Get("1")
	lfu.Get("1")
	lfu.Get("2")

	set_success := lfu.Set(0, "A")
	if !set_success {
		t.Errorf("Failed to set binding with key: %s", "A")
		t.FailNow()
	}

	get_success := lfu.Get("4")
	if get_success {
		t.Errorf("Item with key '4' should have been evicted.")
		t.FailNow()
	}

	set_success2 := lfu.Set(0, "B")
	if !set_success2 {
		t.Errorf("Failed to set binding with key: %s", "B")
		t.FailNow()
	}

	get_success2 := lfu.Get("A")
	if get_success2 {
		t.Errorf("Item with key 'A' should have been evicted.")
		t.FailNow()
	}
}

// Tests whether the LFU cache is working.
func TestLFUFunction(t *testing.T) {
	capacity := 5
	cache := NewLFUCache(capacity)

	var values [5]string
	values[0] = "a"
	values[1] = "b"
	values[2] = "c"
	values[3] = "d"
	values[4] = "e"

	for i := 0; i < 5; i++ {
		key := values[i]
		ok := cache.Set(i, key)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", values[i])
			t.FailNow()
		}
	}

	cache.Get("b")
	cache.Get("a")
	cache.Get("a")
	cache.Get("d")
	cache.Get("e")

	cache.Set(1, "f")

	ok := cache.Get("c")

	if ok {
		t.Errorf("c should have been evicted.")
		t.FailNow()
	}
	cache.Set(1, "g")

	ok = cache.Get("f")

	if ok {
		t.Errorf("f should have been evicted.")
		t.FailNow()
	}

	cache.Get("g")
	cache.Get("g")
	cache.Get("d")
	cache.Get("e")

	cache.Set(1, "h")

	ok = cache.Get("b")

	if ok {
		t.Errorf("b should have been evicted.")
		t.FailNow()
	}
}
