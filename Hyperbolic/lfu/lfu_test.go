/******************************************************************************
 * hyperbolic_test.go
 * Authors: Reuben Agogoe, Stephen Dong, Jimmy Hoang
 * Usage: `go test`  or  `go test -v`
 * Description: A testing suite for comparing hyperbolic caching to FIFO, LRU, and LFU.
 ******************************************************************************/

package cache

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

// // TestHitRate computes the hit miss ratio of the FIFO, LFU, LRU, and Hyperbolic
// // caching algorithms on traces from https://github.[com/twitter/cache-trace.
func TestHitRate(t *testing.T) {

	// traces we will use for testing and evaluation
	traces_to_process := []string{"test.tr"}

	// for i := 2; i < 17; i++ {
	// 	cluster := "cluster0"

	// 	if i < 10 {
	// 		cluster += "0" + strconv.Itoa(i)
	// 	} else {
	// 		cluster += strconv.Itoa(i)
	// 	}

	// 	traces_to_process = append(traces_to_process, cluster)
	// }

	// from the academic paper on Hypbolic caching
	sample_size := 64

	// from Caches Precept
	//max_capacities := []int{80, 160, 320, 640, 1280, 2560, 5120, 10240}
	max_capacities := []int{100}

	for _, max_capacity := range max_capacities {
		for _, trace := range traces_to_process {

			//trace_file := filepath.Join("traces", trace)

			fmt.Println("Testing max capacity: ", max_capacity, " ---")

			RunCacheExperiment(trace, "FIFO", max_capacity, sample_size)
			RunCacheExperiment(trace, "LRU", max_capacity, sample_size)
			RunCacheExperiment(trace, "HYPERBOLIC", max_capacity, sample_size)
			RunCacheExperiment(trace, "LFU", max_capacity, sample_size)

			fmt.Println()
		}
	}
}

// RunCacheExperiment runs an experiment on an input trace file using
// the given cache type, max cache capacity, and (if applicable) sample size.
func RunCacheExperiment(trace_file string, cache_type string, capacity int, sample_size int) {

	// open the trace file
	file, err := os.Open(trace_file)

	if err != nil {
		log.Fatal(err)
	}

	// create a new cache_type cache
	var cache Cache
	cache = NewLFUCache(capacity)

	// read each line of the trace file, parsing the relevant fields
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		text := scanner.Text()

		// format: timestamp, anonymized key, key size,
		// value size, client id, operation, TTL

		line := strings.Split(text, " ")

		_, key, _ :=
			line[0], line[1], line[2]

			// convert string timestamp to int
			//operation_timestamp, _ := strconv.Atoi(timestamp)

			// only handle get and set operations
		_, get_success := cache.Get(key)

		// set if get failed
		if !get_success {
			set_success := cache.Set(key)

			if !set_success {
				log.Fatal("Failed to complete the set request.")
			}
		}
		time.Sleep(2 * time.Millisecond)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// get stats and print them out
	stats := cache.Stats()
	// fmt.Println(stats)
	fmt.Println(cache_type, " Hit Ratio: ",
		float32(stats.Hits)/(float32(stats.Hits+stats.Misses)))
}

//  func TestLFU(t *testing.T) {

// 	 lfu := NewLFUCache(5)

// 	 lfu.Set(1, "a")
// 	 lfu.Set(1, "a")

// 	 lfu.Set(1, "b")

// 	 lfu.Set(1, "c")

// 	 lfu.Set(1, "d")

// 	 lfu.Set(1, "e")

// 	 fmt.Println(lfu.size, "shoul be 5")

// 	 lfu.Get("b")
// 	 lfu.Get("b")

// 	 lfu.Get("c")

// 	 lfu.Get("d")
// 	 lfu.Get("d")
// 	 lfu.Get("d")

// 	 lfu.Get("e")
// 	 lfu.Get("e")
// 	 lfu.Get("e")
// 	 lfu.Get("e")

// 	 lfu.Set(1, "f")
// 	 lfu.Set(1, "g")

// 	 fmt.Println(lfu.size)

// 	 ok := lfu.Get("a")
// 	 ok1 := lfu.Get("b")
// 	 ok2 := lfu.Get("c")
// 	 ok3 := lfu.Get("d")
// 	 ok4 := lfu.Get("e")

// 	 if !ok {
// 		 fmt.Println("A: CORRECT!")
// 	 }
// 	 if !ok1 {
// 		 fmt.Println("B: INCORRECT!")
// 	 }
// 	 if !ok2 {
// 		 fmt.Println("C: CORRECT!")
// 	 }
// 	 if !ok3 {
// 		 fmt.Println("D: INCORRECT!")
// 	 }
// 	 if !ok4 {
// 		 fmt.Println("E: INCORRECT!")
// 	 }

// 	 stats := lfu.Stats()
// 	 fmt.Println(stats)

//  }

//  func TestLFU2(t *testing.T) {

// 	 fmt.Println("STEPHENS CODE -------")
// 	 lfu := NewLFUCache(3)
// 	 fmt.Println(lfu.Get("a"))
// 	 lfu.Set(1, "a")

// 	 fmt.Println(lfu.Get("b"))
// 	 lfu.Set(1, "b")

// 	 fmt.Println(lfu.Get("c"))
// 	 lfu.Set(1, "c")

// 	 fmt.Println(lfu.Get("c"))
// 	 fmt.Println(lfu.Get("c"))
// 	 fmt.Println(lfu.Get("a"))
// 	 lfu.Set(1, "d")
// 	 fmt.Println(lfu.Get("a"))
// 	 fmt.Println(lfu.Get("b"))
//  }
