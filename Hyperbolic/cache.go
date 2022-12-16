package cache

type Stats struct {
	Hits   int
	Misses int
}

type Cache interface {

	// Get returns a success boolean indicating if an item with 
	// the given key was found.
	Get(key string) (ok bool)

	// Set adds/updates an item with the given key in the cache
	// and returns a success boolean.
	Set(operation_timestamp int, key string) (ok bool)

	// Stats returns a pointer to a Stats object that indicates how many hits
	// and misses this cache has resolved over its lifetime.
	Stats() *Stats
}
