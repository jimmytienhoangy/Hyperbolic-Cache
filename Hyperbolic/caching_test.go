/******************************************************************************
* caching_test.go
* Authors: Reuben Agogoe, Stephen Dong, Jimmy Hoang
* Usage: `go test`  or  `go test -v`
* Description: A testing suite for comparing hyperbolic caching to the
				FIFO, LRU, and LFU caching algorithms with traces.
******************************************************************************/

package cache

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestHitRate computes the hit ratio of the FIFO, LFU, LRU, and Hyperbolic
// caching algorithms on traces from https://github.com/twitter/cache-trace.
func TestHitRate(t *testing.T) {

	// traces we will use for testing and evaluation
	traces_to_process := []string{"cluster052"}

	// constant given by the academic paper on hyperbolic caching
	// linked in the final project document on the course website
	sample_size := 64

	// max capacities to test and evaluate
	max_capacities := []int{10000}

	// run each caching algorithm with every combination of input
	// trace files and max capacities
	for _, max_capacity := range max_capacities {
		for _, trace := range traces_to_process {

			trace_file := filepath.Join("traces", trace)

			fmt.Println("Testing max capacity [", max_capacity, "] on "+
				"trace file ["+trace_file+"] ---")

			RunCacheExperiment(trace_file, "FIFO", max_capacity, sample_size)
			RunCacheExperiment(trace_file, "LRU", max_capacity, sample_size)
			RunCacheExperiment(trace_file, "LFU", max_capacity, sample_size)
			RunCacheExperiment(trace_file, "HYPERBOLIC", max_capacity, sample_size)

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

	if cache_type == "FIFO" {
		cache = NewFIFOCache(capacity)
	} else if cache_type == "LRU" {
		cache = NewLRUCache(capacity)
	} else if cache_type == "LFU" {
		cache = NewLFUCache(capacity)
	} else if cache_type == "HYPERBOLIC" {
		cache = NewHyperbolicCache(capacity, sample_size)
	}

	// read each line of the trace file, parsing the relevant fields
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		text := scanner.Text()

		// format: timestamp, anonymized key, key size,
		// value size, client id, operation, TTL

		line := strings.Split(text, ",")

		timestamp, key, _, _, _, operation, _ :=
			line[0], line[1], line[2], line[3], line[4], line[5], line[6]

		// convert string timestamp to int
		operation_timestamp, _ := strconv.Atoi(timestamp)

		// only handle get and set operations
		if operation == "set" {
			set_success := cache.Set(operation_timestamp, key)

			if !set_success {
				log.Fatal("Failed to complete the set request.")
			}
		} else if operation == "get" {
			get_success := cache.Get(key)

			// set if get failed
			if !get_success {
				set_success := cache.Set(operation_timestamp, key)

				if !set_success {
					log.Fatal("Failed to complete the set request.")
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// get stats and print them out
	stats := cache.Stats()

	fmt.Println(cache_type, "Hit Ratio:",
		float32(stats.Hits)/(float32(stats.Hits+stats.Misses)))
}
