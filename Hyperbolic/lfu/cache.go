package cache

type Stats struct {
	Hits   int
	Misses int
}

type Cache interface {
	// MaxStorage returns the maximum number of bytes this cache can store.
	MaxStorage() int

	// RemainingStorage returns the number of unused bytes available in this cache.
	RemainingStorage() int

	// Get returns the value_size associated with the given key, if it exists.
	// This operation counts as a "use" for the item with the given key.
	// ok is true if a value_size was found and false otherwise.
	Get(key string) (access_count int, ok bool)

	// Remove removes the key-value_size pair associated with the given key, if it exists.
	// ok is true if a value_size was found and false otherwise.
	Remove(key string) (ok bool)

	// Set associates the given value_size with the given key, possibly evicting key-value_size pairs
	// to make room. This operation counts as a "use" for the item with the given key.
	// Returns true if the binding was added successfully, else false.
	Set(key string) (ok bool)

	// Len returns the number of bindings in the cache.
	Len() int

	// Stats returns a pointer to a Stats object that indicates how many hits
	// and misses this cache has resolved over its lifetime.
	Stats() *Stats
}
