/******************************************************************************
 * fifo_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for fifo.go. You are welcome to change
 *    anything in this file however you would like. You are strongly encouraged
 *    to create additional tests for your implementation, as the ones provided
 *    here are extremely basic, and intended only to demonstrate how to test
 *    your program.
 ******************************************************************************/

package cache

import (
	"fmt"
	"testing"
)

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/
// Constants can go here

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

// Test to see if a normal FIFO cache that doesn't reach capacity is working
// (get and set methods)
func TestFIFO1(t *testing.T) {
	capacity := 64
	fifo := NewFifo(capacity)
	checkCapacity(t, fifo, capacity)

	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := fifo.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}

		res, _ := fifo.Get(key)
		if !bytesEqual(res, val) {
			t.Errorf("Wrong value %s for binding with key: %s", res, key)
			t.FailNow()
		}
	}
}

// Check to see that eviction is occuring at all
func TestFIFO2(t *testing.T) {
	capacity := 64
	fifo := NewFifo(capacity)
	checkCapacity(t, fifo, capacity)

	for i := 0; i < 8; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := fifo.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}

	}

	// here is where the eviction is happening
	for i := 9; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := fifo.Set(key, val)
		if !ok {
			t.Errorf("Failed to evict in order to make space for new caces %s", key)
			t.FailNow()
		}
	}
}

// Test a 0 capacity FIFO
func TestFIFO3(t *testing.T) {
	capacity := 0
	fifo := NewFifo(capacity)
	checkCapacity(t, fifo, capacity)

	// try to set bindings
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := fifo.Set(key, val)

		if ok {
			t.Errorf("There should not be any bindings inside the FIFO!")
			t.FailNow()
		}
	}
}

// Check to see if a new FIFO returns a fifo of empty size
func TestFIFO4(t *testing.T) {
	capacity := 16
	fifo := NewFifo(capacity)
	checkCapacity(t, fifo, capacity)

	// check if new fifo is empty size
	if fifo.Len() != 0 {
		t.Errorf("New fifo cache does not initialize to empty size!")
		t.FailNow()
	}
}

// Test a lot of stuff!
func TestFIFO5(t *testing.T) {
	capacity := 100
	fifo := NewFifo(capacity)
	checkCapacity(t, fifo, capacity)

	// test states and cities to insert as bindings
	var test_states [5]string
	var test_cities [5]string

	test_states[0] = "Alabama"     // 7
	test_states[1] = "Louisiana"   // 9
	test_states[2] = "Mississippi" // 11
	test_states[3] = "Florida"     // 7
	test_states[4] = "Tennessee"   // 9

	test_cities[0] = "Mobile"         // 6
	test_cities[1] = "New Orleans"    // 11
	test_cities[2] = "Pass Christian" // 14
	test_cities[3] = "Orlando"        // 7
	test_cities[4] = "Chattanooga"    // 11

	// set bindings and check remaining storage at some steps
	for i := 0; i < 5; i++ {
		val := []byte(test_cities[i])
		ok := fifo.Set(test_states[i], val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", test_states[i])
			t.FailNow()
		}
		if i == 0 && fifo.RemainingStorage() != 87 {
			t.Errorf("Remaining storage should be 87.")
			t.FailNow()
		}
		if i == 4 && fifo.RemainingStorage() != 8 {
			t.Errorf("Remaining storage should be 8.")
			t.FailNow()
		}
	}

	// remove a binding and check storage and binding count
	fifo.Remove("Alabama")
	if fifo.RemainingStorage() != 21 {
		t.Errorf("Remaining storage should be 21.")
		t.FailNow()
	}
	if fifo.Len() != 4 {
		t.Errorf("There should be 4 bindings.")
		t.FailNow()
	}

	// update the value for an existing key
	fifo.Set("Mississippi", []byte("Biloxi"))

	// check number of bindings and storage
	if fifo.Len() != 4 {
		t.Errorf("There should still be 4 bindings.")
		t.FailNow()
	}
	if fifo.RemainingStorage() != 29 {
		t.Errorf("Remaining storage should be 29.")
		t.FailNow()
	}
}
