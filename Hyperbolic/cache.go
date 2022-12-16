package cache

type Stats struct {
	Hits   int
	Misses int
}

type Cache interface {

	// Get returns true if an item with the given
	// key was found in the cache, false otherwise.
	Get(key string) (success bool)

	// Set adds or updates an item with the given key in the
	// cache and returns true if a successful update was
	// made, false otherwise.
	Set(operation_timestamp int, key string) (success bool)

	// Stats returns a pointer to a Stats object that indicates how many hits
	// and misses this cache has resolved over its lifetime.
	Stats() *Stats
}
