package cache

type Stats struct {
	Hits   int
	Misses int
}

func (stats *Stats) Equals(other *Stats) bool {
	if stats == nil && other == nil {
		return true
	}
	if stats == nil || other == nil {
		return false
	}
	return stats.Hits == other.Hits && stats.Misses == other.Misses
}

type Cache interface {
	// MaxStorage returns the maximum number of bytes this cache can store
	MaxStorage() int

	// RemainingStorage returns the number of unused bytes available in this cache
	RemainingStorage() int

	// Get returns the value associated with the given key, if it exists.
	// This operation counts as a "use" for that key-value pair
	// ok is true if a value was found and false otherwise.
	Get(key string) (value []byte, ok bool)

	// Remove removes and returns the value associated with the given key, if it exists.
	// ok is true if a value was found and false otherwise
	Remove(key string) (value []byte, ok bool)

	// Set associates the given value with the given key, possibly evicting values
	// to make room. Returns true if the binding was added successfully, else false.
	Set(key string, value []byte) bool

	// Len returns the number of bindings in the cache.
	Len() int

	// Stats returns a pointer to a Stats object that indicates how many hits
	// and misses this cache has resolved over its lifetime.
	Stats() *Stats
}
