package cache

import (
	"testing"
)

/******************************************************************************/
/*                                 Helpers                                    */
/******************************************************************************/

// CacheType returns a string representing the type (i.e. eviction scheme) of
// this cache.
func cacheType(cache Cache) string {
	switch cache.(type) {
	case *LRU:
		return "LRU"
	case *FIFO:
		return "FIFO"
	default:
		return "cache"
	}
}

// Returns true iff a and b represent equal slices of bytes.
func bytesEqual(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}

	return true
}

// Fails test t with an error message if fifo.MaxStorage() is not equal to capacity
func checkCapacity(t *testing.T, cache Cache, capacity int) {
	max := cache.MaxStorage()
	if max != capacity {
		t.Errorf("Expected %s to have %d MaxStorage, but it had %d", cacheType(cache), capacity, max)
	}
}
