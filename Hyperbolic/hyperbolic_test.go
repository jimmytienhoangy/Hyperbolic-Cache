/******************************************************************************
 * hyperbolic_test.go
 * Authors: Reuben Agogoe, Stephen Dong, Jimmy Hoang
 * Usage: `go test`  or  `go test -v`
 * Description: A unit testing suite for hyperbolic.go.
 ******************************************************************************/

package cache

// import (
// 	"fmt"
// 	"testing"
// )

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

// // Test the creation of a hyperbolic cache
// func Test_CreateHyperbolic(t *testing.T) {
// 	capacity := 64
// 	hyperbolic := NewHyperbolicCache(capacity, 64)
// 	for i := 0; i < 4; i++ {
// 		key := fmt.Sprintf("key%d", i)
// 		ok := hyperbolic.Set(0, key)
// 		if !ok {
// 			t.Errorf("Failed to add binding with key: %s", key)
// 			t.FailNow()
// 		}
// 		ok = hyperbolic.Get(key)
// 	}
// }

// // Check to see that eviction is occuring at all to make spaces for out of capacity
// func Test_HyperbolicEviction(t *testing.T) {
// 	capacity := 3
// 	hyperbolic := NewHyperbolicCache(capacity, 3)

// 	for i := 0; i < 3; i++ {
// 		key := fmt.Sprintf("key%d", i)

// 		ok := hyperbolic.Set(i, key)
// 		if !ok {
// 			t.Errorf("Failed to add binding with key: %s", key)
// 			t.FailNow()
// 		}

// 	}
// 	ok := hyperbolic.Set(5, "evict_for_space")
// 	if !ok {
// 		t.Errorf("Failed to add binding with key: %s", "evict_for_space")
// 		t.FailNow()
// 	}
// }

// // // Test a 0 capacity hyperbolic
// // func TestHyperbolic3(t *testing.T) {
// // 	capacity := 0
// // 	hyperbolic := HyperbolicCache(capacity, 0)

// // 	// try to set bindings
// // 	for i := 0; i < 3; i++ {
// // 		key := fmt.Sprintf("key%d", i)
// // 		val := []byte(key)
// // 		ok := lru.Set(key, val)

// // 		if ok {
// // 			t.Errorf("There should not be any bindings inside the lru!")
// // 			t.FailNow()
// // 		}
// // 	}
// // }

// // // Check to see if a new lru returns a lru of empty size
// // func TestHyperbolic4(t *testing.T) {
// // 	capacity := 16
// // 	lru := NewLru(capacity)
// // 	checkCapacity(t, lru, capacity)

// // 	// check if new lru is empty size
// // 	if lru.Len() != 0 {
// // 		t.Errorf("New lru cache does not initialize to empty size!")
// // 		t.FailNow()
// // 	}
// // }

// // sleep for just 1 ms lol
// /*func wait() {
// 	time.Sleep(100 * time.Millisecond)
// }*/

// // Test whether the hyperbolic function is working
// func TestHyperbolicFunction(t *testing.T) {
// 	capacity := 5
// 	cache := NewHyperbolicCache(capacity, 5)
// 	// test states and cities to insert as bindings
// 	var values [5]string
// 	values[0] = "a" // 7
// 	values[1] = "b" // 9
// 	values[2] = "c" // 11
// 	values[3] = "d" // 7
// 	values[4] = "e" // 9
// 	// set bindings and check remaining storage at some steps
// 	for i := 0; i < 5; i++ {
// 		key := values[i]
// 		timestamp := (i + 1) * 2
// 		ok := cache.Set(timestamp, key)
// 		if !ok {
// 			t.Errorf("Failed to add binding with key: %s", values[i])
// 			t.FailNow()
// 		}
// 	}

// 	cache.Get("b")
// 	cache.Get("a")
// 	cache.Get("a")

// 	cache.Set((5+1)*2, "f")

// 	for j := range cache.keys_to_items {
// 		fmt.Println(j)
// 	}

// 	ok := cache.Get("c")

// 	if ok {
// 		t.Errorf("c should have been evicted.")
// 		t.FailNow()
// 	}
// 	fmt.Println(cache.Stats())
// }


// // Test whether the hyperbolic function is working (no set)
// func TestHyperbolicFunction2(t *testing.T) {
// 	capacity := 5
// 	cache := NewHyperbolicCache(capacity, 5)
// 	// test states and cities to insert as bindings
// 	var values [5]string
// 	values[0] = "a" // 7
// 	values[1] = "b" // 9
// 	values[2] = "c" // 11
// 	values[3] = "d" // 7
// 	values[4] = "e" // 9
// 	// set bindings and check remaining storage at some steps
// 	for i := 0; i < 5; i++ {
// 		key := values[i]
// 		timestamp := (i + 1) * 2
// 		ok := cache.Set(timestamp, key)
// 		if !ok {
// 			t.Errorf("Failed to add binding with key: %s", values[i])
// 			t.FailNow()
// 		}
// 	}

// 	cache.Set((5+1)*2, "f")

// 	for j := range cache.keys_to_items {
// 		fmt.Println(j)
// 	}

// 	ok := cache.Get("a")

// 	if ok {
// 		t.Errorf("a should have been evicted.")
// 		t.FailNow()
// 	}
// 	fmt.Println(cache.Stats())
// }
// /******************************************************************************
//  * hyperbolic_test.go
//  * Authors: Reuben Agogoe, Stephen Dong, Jimmy Hoang
//  * Usage: `go test`  or  `go test -v`
//  * Description: A testing suite for comparing hyperbolic caching to 
//  				FIFO, LRU, and LFU.
//  ******************************************************************************/

// package cache

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"
// 	"testing"
// )

// // // TestHitRate computes the hit miss ratio of the FIFO, LFU, LRU, and Hyperbolic
// // // caching algorithms on traces from https://github.[com/twitter/cache-trace.
// func TestHitRate(t *testing.T) {

// 	// traces we will use for testing and evaluation
// 	traces_to_process := []string{"cluster052"}

// 	// for i := 2; i < 17; i++ {
// 	// 	cluster := "cluster0"

// 	// 	if i < 10 {
// 	// 		cluster += "0" + strconv.Itoa(i)
// 	// 	} else {
// 	// 		cluster += strconv.Itoa(i)
// 	// 	}

// 	// 	traces_to_process = append(traces_to_process, cluster)
// 	// }

// 	// from the academic paper on Hypbolic caching
// 	sample_size := 64

// 	// from Caches Precept
// 	max_capacities := []int{100, 200, 300, 400, 500, 1000, 2000, 3000, 4000}

// 	for _, max_capacity := range max_capacities {
// 		for _, trace := range traces_to_process {

// 			trace_file := filepath.Join("traces", trace)

// 			fmt.Println("Testing max capacity: ", max_capacity, " ---")

// 			RunCacheExperiment(trace_file, "FIFO", max_capacity, sample_size)
// 			RunCacheExperiment(trace_file, "LRU", max_capacity, sample_size)
// 			RunCacheExperiment(trace_file, "HYPERBOLIC", max_capacity, sample_size)
// 			RunCacheExperiment(trace_file, "LFU", max_capacity, sample_size)

// 			fmt.Println()
// 		}
// 	}
// }

// // RunCacheExperiment runs an experiment on an input trace file using
// // the given cache type, max cache capacity, and (if applicable) sample size.
// func RunCacheExperiment(trace_file string, cache_type string, capacity int, sample_size int) {

// 	// open the trace file
// 	file, err := os.Open(trace_file)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// create a new cache_type cache
// 	var cache Cache

// 	if cache_type == "FIFO" {
// 		cache = NewFIFOCache(capacity)
// 	} else if cache_type == "LRU" {
// 		cache = NewLRUCache(capacity)
// 	} else if cache_type == "LFU" {
// 		cache = NewLFUCache(capacity)
// 	} else if cache_type == "HYPERBOLIC" {
// 		cache = NewHyperbolicCache(capacity, sample_size)
// 	}

// 	// read each line of the trace file, parsing the relevant fields
// 	defer file.Close()
// 	scanner := bufio.NewScanner(file)

// 	for scanner.Scan() {

// 		text := scanner.Text()

// 		// format: timestamp, anonymized key, key size,
// 		// value size, client id, operation, TTL

// 		line := strings.Split(text, ",")

// 		timestamp, key, _, _, _, operation, _ :=
// 			line[0], line[1], line[2], line[3], line[4], line[5], line[6]

// 		// convert string timestamp to int
// 		operation_timestamp, _ := strconv.Atoi(timestamp)

// 		// only handle get and set operations
// 		if operation == "set" {
// 			set_success := cache.Set(operation_timestamp, key)

// 			if !set_success {
// 				log.Fatal("Failed to complete the set request.")
// 			}
// 		} else if operation == "get" {
// 			get_success := cache.Get(key)

// 			// set if get failed
// 			if !get_success {
// 				set_success := cache.Set(operation_timestamp, key)

// 				if !set_success {
// 					log.Fatal("Failed to complete the set request.")
// 				}
// 			}
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	// get stats and print them out
// 	stats := cache.Stats()
// 	// fmt.Println(stats)
// 	fmt.Println(cache_type, " Hit Ratio: ",
// 		float32(stats.Hits)/(float32(stats.Hits+stats.Misses)))

// }