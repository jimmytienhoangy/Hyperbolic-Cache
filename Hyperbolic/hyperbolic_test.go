/******************************************************************************
 * hyperbolic_test.go
 * Authors: Reuben Agogoe, Stephen Dong, Jimmy Hoang
 * Usage: `go test`  or  `go test -v`
 * Description: A unit testing suite for hyperbolic.go.
 ******************************************************************************/

package cache

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"strings"
	"testing"
	//"time"
)

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

// Test to see if a normal LRU cache that doesn't reach capacity is working
// (get and set methods)
/*func TestHyperbolic(t *testing.T) {
	capacity := 64
	hyperbolic := NewHyperbolicCache(capacity)

	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("key%d", i)
		val := len([]byte(key))

		ok := hyperbolic.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}

		res, _ := hyperbolic.Get(key)
		// if !bytesEqual(res, val) {
		// 	t.Errorf("Wrong value %s for binding with key: %s", res, key)
		// 	t.FailNow()
		// }

		if res != val {
			t.Errorf("Wrong value %d for binding with key: %s", res, key)
			t.FailNow()

		}
	}
}*/

// // Check to see that eviction is occuring at all
// func TestHyperbolic2(t *testing.T) {
// 	capacity := 64
// 	lru := NewLru(capacity)
// 	checkCapacity(t, lru, capacity)

// 	for i := 0; i < 8; i++ {
// 		key := fmt.Sprintf("key%d", i)
// 		val := []byte(key)
// 		ok := lru.Set(key, val)
// 		if !ok {
// 			t.Errorf("Failed to add binding with key: %s", key)
// 			t.FailNow()
// 		}

// 	}

// 	// here is where the eviction is happening
// 	for i := 9; i < 10; i++ {
// 		key := fmt.Sprintf("key%d", i)
// 		val := []byte(key)
// 		ok := lru.Set(key, val)
// 		if !ok {
// 			t.Errorf("Failed to evict in order to make space for new caces %s", key)
// 			t.FailNow()
// 		}
// 	}
// }

// // Test a 0 capacity lru
// func TestHyperbolic3(t *testing.T) {
// 	capacity := 0
// 	lru := (capacity)
// 	checkCapacity(t, lru, capacity)

// 	// try to set bindings
// 	for i := 0; i < 3; i++ {
// 		key := fmt.Sprintf("key%d", i)
// 		val := []byte(key)
// 		ok := lru.Set(key, val)

// 		if ok {
// 			t.Errorf("There should not be any bindings inside the lru!")
// 			t.FailNow()
// 		}
// 	}
// }

// // Check to see if a new lru returns a lru of empty size
// func TestHyperbolic4(t *testing.T) {
// 	capacity := 16
// 	lru := NewLru(capacity)
// 	checkCapacity(t, lru, capacity)

// 	// check if new lru is empty size
// 	if lru.Len() != 0 {
// 		t.Errorf("New lru cache does not initialize to empty size!")
// 		t.FailNow()
// 	}
// }

// Test a lot of stuff!
// func TestHyperbolic5(t *testing.T) {
// 	capacity := 100
// 	lru := NewHyperbolicCache(capacity)
// 	checkCapacity(t, lru, capacity)

// 	// test states and cities to insert as bindings
// 	var test_states [5]string
// 	var test_cities [5]string

// 	test_states[0] = "Alabama"     // 7
// 	test_states[1] = "Louisiana"   // 9
// 	test_states[2] = "Mississippi" // 11
// 	test_states[3] = "Florida"     // 7
// 	test_states[4] = "Tennessee"   // 9

// 	test_cities[0] = "Mobile"         // 6
// 	test_cities[1] = "New Orleans"    // 11
// 	test_cities[2] = "Pass Christian" // 14
// 	test_cities[3] = "Orlando"        // 7
// 	test_cities[4] = "Chattanooga"    // 11

// 	// set bindings and check remaining storage at some steps
// 	for i := 0; i < 5; i++ {
// 		val := []byte(test_cities[i])
// 		ok := lru.Set(test_states[i], val)
// 		if !ok {
// 			t.Errorf("Failed to add binding with key: %s", test_states[i])
// 			t.FailNow()
// 		}
// 		if i == 0 && lru.RemainingStorage() != 87 {
// 			t.Errorf("Remaining storage should be 87.")
// 			t.FailNow()
// 		}
// 		if i == 4 && lru.RemainingStorage() != 8 {
// 			t.Errorf("Remaining storage should be 8.")
// 			t.FailNow()
// 		}
// 	}

// 	// remove a binding and check storage and binding count
// 	lru.Remove("Alabama")
// 	if lru.RemainingStorage() != 21 {
// 		t.Errorf("Remaining storage should be 21.")
// 		t.FailNow()
// 	}
// 	// if lru.Len() != 4 {
// 	// 	t.Errorf("There should be 4 bindings.")
// 	// // 	t.FailNow()
// 	// }

// 	// update the value for an existing key
// 	lru.Set("Mississippi", []byte("Biloxi"))

// 	// // check number of bindings and storage
// 	// if lru.Len() != 4 {
// 	// 	t.Errorf("There should still be 4 bindings.")
// 	// 	t.FailNow()
// 	// }
// 	if lru.RemainingStorage() != 29 {
// 		t.Errorf("Remaining storage should be 29.")
// 		t.FailNow()
// 	}

// 	// test eviction
// 	lru.Set("California", []byte("San Francisco")) // 23
// 	lru.Set("Massachusetts", []byte("Boston"))     // 19

// 	_, another_ok := lru.Get("Louisiana")

// 	if another_ok {
// 		t.Errorf("Louisiana should have been evicted.")
// 		t.FailNow()
// 	}

// 	fmt.Println(lru.Stats())

// }

// sleep for just 1 ms lol
/*func wait() {
	time.Sleep(100 * time.Millisecond)
}*/

// Test whether the hyperbolic function is working
/*func TestHyperbolicFunction(t *testing.T) {
	capacity := 10
	cache := NewHyperbolicCache(capacity)

	// test states and cities to insert as bindings
	var values [5]string

	values[0] = "a" // 7
	values[1] = "b" // 9
	values[2] = "c" // 11
	values[3] = "d" // 7
	values[4] = "e" // 9

	// set bindings and check remaining storage at some steps
	for i := 0; i < 5; i++ {
		key := values[i]
		val := len([]byte(key))
		ok := cache.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", values[i])
			t.FailNow()
		}
		wait()

	}
	fmt.Println("test")
	cache.Get("b")
	wait()
	cache.Get("a")
	wait()
	cache.Get("a")
	wait()
	fmt.Println("test2")

	cache.Set("f", len([]byte("f")))
	fmt.Println("test2")

	for j := range cache.keys_to_items {
		fmt.Println(j)
	}

	_, ok := cache.Get("c")

	if ok {
		t.Errorf("c should have been evicted.")
		t.FailNow()
	}

	fmt.Println(cache.Stats())

}

// Test whether the hyperbolic function is working
func TestHyperbolicFunction2(t *testing.T) {
	capacity := 10
	cache := NewHyperbolicCache(capacity)

	// test states and cities to insert as bindings
	var values [3]string

	values[0] = "a"
	values[1] = "bc"
	values[2] = "cd"

	// set bindings and check remaining storage at some steps
	for i := 0; i < 3; i++ {
		key := values[i]
		val := len([]byte(key))
		ok := cache.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", values[i])
			t.FailNow()
		}
		if i == 0 {
			_, ok := cache.Get("b")

			if ok {
				fmt.Println("Error: item with key 'b' found.")
			}
		}
		wait()
	}

	cache.Get("bc")
	wait()

	cache.Get("bc")
	wait()

	cache.Get("bc")
	wait()

	cache.Get("bc")
	wait()

	cache.Get("bc")
	wait()
	cache.Set("de", len([]byte("de")))
	wait()

	_, ok1 := cache.Get("a")

	if ok1 {
		t.Errorf("Item with key 'a' should have been evicted.")
		t.FailNow()
	}

	_, ok2 := cache.Get("cd")

	for j := range cache.keys_to_items {
		fmt.Println(j)
	}

	if ok2 {
		t.Errorf("Item with key 'cd' should have been evicted.")
		t.FailNow()
	}

	fmt.Println(cache.Stats())

}*/

// func TestHyperbolicFunction3(t *testing.T) {

// }

func simulateHyperbolicCache(t *testing.T) {

	// read in trace file
	path := filepath.Join("cache-trace", "samples", "2020Mar", "cluster002")
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	// create a hyperbolic cache with max_capacity, sample size
	max_capacity := 10000
	sample_size := 64
	hyperbolic_cache := NewHyperbolicCache(max_capacity, sample_size)

	defer file.Close()
	scanner := bufio.NewScanner(file)

	// read each line of the trace file, parse relevant fields, and fulfill request
	for scanner.Scan() {

		text := scanner.Text()

		line := strings.Split(text, ",")

		_, key, _, _, _, operation, _ :=
			line[0], line[1], line[2], line[3], line[4], line[5], line[6]

		// only handle get and set operations
		if operation == "set" {
			set_success := hyperbolic_cache.Set(key)
			if !set_success {
				log.Fatal("Failed to complete the set request.")
			}
		} else if operation == "get" {
			// we can return access count if we want to check accuracy
			_, get_success := hyperbolic_cache.Get(key)

			if !get_success {
				set_success := hyperbolic_cache.Set(key)

				if !set_success {
					log.Fatal("Failed to complete the set request.")
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// print hit, miss ratio
	fmt.Println(hyperbolic_cache.Stats())
}
