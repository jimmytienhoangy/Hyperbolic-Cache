package cache

type Stats struct {
	Hits   int
	Misses int
}

type Cache interface {
	// MaxStorage returns the maximum number of items this cache can store.
	MaxStorage() int

	// RemainingStorage returns the number of items that can still be stored
	// in this cache.
	RemainingStorage() int

	// Get returns how many times the item associated with the given key
	// has been accessed (not including this access) and a success boolean.
	Get(key string) (access_count int, ok bool)

	// Remove removes the item associated with the given key from the cache, if it exists.
	// ok is true if an item was found and false otherwise.
	Remove(key string) (ok bool)

	// Set adds/updates an item with the given key in the cache
	// and returns a success boolean. This operation counts as a
	// access for the item with the given key.
	Set(key string) (ok bool)

	// Len returns the number of items in the cache.
	Len() int

	// Stats returns a pointer to a Stats object that indicates how many hits
	// and misses this cache has resolved over its lifetime.
	Stats() *Stats
}
