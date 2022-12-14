package cache

import (
	"testing"
)

/******************************************************************************/
/*                                 Helpers                                    */
/******************************************************************************/

// Fails test t with an error message if hyperbolic.MaxStorage() is not equal to capacity
func checkCapacity(t *testing.T, cache Cache, capacity int) {
	max_storage := cache.MaxStorage()
	if max_storage != capacity {
		t.Errorf("Expected the cache to have %d max capacity, but it had %d!", capacity, max_storage)
	}
}
