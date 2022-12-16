package cache

import (
	"container/list"
)

// A LFUCacheItem is an item with metadata that implicitly
// holds a value of set size. It goes in a LFUCache.
type LFUCacheItem struct {

	// the item's key
	key string

	// pointer to a struct indicating this item's access count
	accessParent *list.Element
}

// An AccessNode is the connection between an access count and
// the items with that access count.
type AccessNode struct {

	// items with this access item's access count
	items_with_access_count map[*LFUCacheItem]byte

	// the access count associated with this access item
	access_count int
}

// A LFUCache is a cache that uses lfu caching.
type LFUCache struct {

	// maximum number of items the cache can hold
	max_capacity int

	// total number of items currently in the cache
	size int

	// map of keys to items in the cache
	keys_to_items map[string]*LFUCacheItem

	// linked list of access counts
	access_counts *list.List

	// number of hits
	hits int

	// number of misses
	misses int
}

// NewLFUCache creates a new, empty LFUCache.
func NewLFUCache(max_capacity int) *LFUCache {

	return &LFUCache{
		max_capacity:  max_capacity,
		size:          0,
		keys_to_items: make(map[string]*LFUCacheItem, max_capacity),
		access_counts: list.New(),
		hits:          0,
		misses:        0,
	}
}

// Get returns a success boolean indicating if an item with the key was found.
func (cache *LFUCache) Get(key string) (success bool) {

	// retrieve item associated with key
	item, ok := cache.keys_to_items[key]

	if ok {
		cache.hits += 1

		// update access count of item
		cache.increment(item)

	} else {
		cache.misses += 1

		return false
	}

	return true
}

// Set adds/updates an item with the given key in the cache
// and returns a success boolean.
func (lfu *LFUCache) Set(operation_timestamp int, key string) (success bool) {

	// can not set if cache max capacity is 0!
	if lfu.max_capacity == 0 {
		return false
	}

	// operation_timestamp is ignored for the purposes of this project

	// check if an item with that key already exists
	existing_item, ok := lfu.keys_to_items[key]

	if ok {
		// update access count of item
		lfu.increment(existing_item)

		return true
	}

	// if not enough space and an item with the key does not exist,
	// evict an item
	if lfu.size == lfu.max_capacity {
		lfu.evict()
	}

	// add new item with key
	new_item := &LFUCacheItem{key: key}
	lfu.keys_to_items[key] = new_item

	// update access for the new item
	lfu.increment(new_item)

	// update size of cache
	lfu.size += 1

	return true
}

// Increment updates the access count of a given item.
func (lfu *LFUCache) increment(item *LFUCacheItem) {

	// check if the item is already associated with an access count node
	currentAccessNode := item.accessParent

	// find new access count value and corresponding node
	var nextAccessCount int
	var nextAccessNode *list.Element

	// item does not have an access count node (new insert!)
	if currentAccessNode == nil {
		// first access
		nextAccessCount = 1
		// item's access count node should be the very front (1st!)
		nextAccessNode = lfu.access_counts.Front()
	} else {
		// increment access count by 1
		nextAccessCount = currentAccessNode.Value.(*AccessNode).access_count + 1

		// move to next access count node (that may not exist)
		nextAccessNode = currentAccessNode.Next()
	}

	// next access count node does not exist or there is a gap:
	// for example, a key has 6 accesses and another with 8 accesses, but no key has 7 accesses anymore
	if nextAccessNode == nil || nextAccessNode.Value.(*AccessNode).access_count != nextAccessCount {

		// create a new access count node for the missing access count
		newAccessNode := new(AccessNode)
		newAccessNode.access_count = nextAccessCount
		newAccessNode.items_with_access_count = make(map[*LFUCacheItem]byte)

		// add new access count node to the front if this is a new insert
		if currentAccessNode == nil {
			nextAccessNode = lfu.access_counts.PushFront(newAccessNode)

		} else {
			// add new access count node after the current one
			nextAccessNode = lfu.access_counts.InsertAfter(newAccessNode, currentAccessNode)
		}
	}

	// set the new access count parent for the item that is being incremented and
	// add it to that parent's list of entries
	item.accessParent = nextAccessNode
	nextAccessNode.Value.(*AccessNode).items_with_access_count[item] = 1

	// remove the item from the entries of its old access count node (currentAccessNode)
	if currentAccessNode != nil {
		lfu.remove(currentAccessNode, item)
	}
}

// evict evicts the least frequently used item from the cache.
func (lfu *LFUCache) evict() {

	// get the smallest access count node
	if smallestAccessNode := lfu.access_counts.Front(); smallestAccessNode != nil {

		// for all the entries of this access count node
		for entry := range smallestAccessNode.Value.(*AccessNode).items_with_access_count {

			// delete the item from the cache
			delete(lfu.keys_to_items, entry.key)

			// remove the item from all lists
			lfu.remove(smallestAccessNode, entry)

			lfu.size--
		}
	}
}

// Remove removes the item associated with the given key from the cache, if it exists.
func (lfu *LFUCache) remove(listItem *list.Element, item *LFUCacheItem) {

	accessNode := listItem.Value.(*AccessNode)

	// remove the item from its corresponding access count node's entries
	delete(accessNode.items_with_access_count, item)

	// this access node no longer has entries, so remove it from the list
	// of access counts
	if len(accessNode.items_with_access_count) == 0 {
		lfu.access_counts.Remove(listItem)
	}
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lfu *LFUCache) Stats() *Stats {
	return &Stats{Hits: lfu.hits, Misses: lfu.misses}
}
