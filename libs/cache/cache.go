// Universal cache implementation using generics to define key and value types
// Can be used with or without eviction
// MaxSize configuration parameter = 0 means no upper size bound (eviction only by TTL if set)
// Eviction can be based on LRU or LFU algorithms
// Eviction on TTL in seconds can be used in combination with any eviction algorithm (TTL eviction triggered only on records addition)
// Cache configuration can be changed on the fly using SetConfiguration method
// All data needed for different eviction methods is always computed and stored in Cache struct, so new configuration starts to work without any troubles
//
// All eviction methods implemented using double linked lists
// Cache records are encapsulated in Node structs
// Node structs are linked in three lists for sorting: LFU, LRU and TTL
// Eviction candidates are nearest to Tail of lists for each method:
//   For LFU LFUCounter incremented on each key use and node position in LFU list is updated according to this counter
//   For LRU node always moved to the head of LRU list on each key use
//   For TTL UsedAt field updated to time.now() and node always moved to the head of TTL list only when Set method used
// When MaxSize > 0 in Config and size of cache is >= MaxSize, then eviction started with TTL if it is set
// LFU or LRU eviction is applied if there is still no space after TTL eviction
//
// Two special methods can be used to invalidate discrete or all records in cache: Invalidate and Clear

package cache

import (
	"context"
	"sync"
	"time"
)

type Node[KeyT comparable, ValueT any] struct { // Node used in cache for value storage and lists linking
	Key      KeyT
	Value    ValueT              // Cached value
	LFUCount uint64              // Counter for LFU eviction
	LFUNext  *Node[KeyT, ValueT] // Next node in double linked list for LFU
	LFUPrev  *Node[KeyT, ValueT] // Prev node in double linked list for LFU
	LRUNext  *Node[KeyT, ValueT] // Next node in double linked list for LRU
	LRUPrev  *Node[KeyT, ValueT] // Prev node in double linked list for LRU
	UsedAt   time.Time           // Time of last node usage
	TTLNext  *Node[KeyT, ValueT] // Next node in double linked list for TTL
	TTLPrev  *Node[KeyT, ValueT] // Prev node in double linked list for TTL
}

type Cache[KeyT comparable, ValueT any] interface {
	Set(ctx context.Context, key KeyT, value ValueT) bool
	Get(ctx context.Context, key KeyT) (*ValueT, bool)
	Invalidate(ctx context.Context, key KeyT) bool
	SetConfig(ctx context.Context, config Config) error
	Clear(ctx context.Context) error
}

type cache[KeyT comparable, ValueT any] struct {
	Lock       sync.RWMutex
	ConfigLock sync.RWMutex
	Config     Config                       // cache configuration parameters
	Size       uint64                       // Current cache size
	Storage    map[KeyT]*Node[KeyT, ValueT] // Index map
	LFUHead    *Node[KeyT, ValueT]          // Head of double linked list for LFU
	LFUTail    *Node[KeyT, ValueT]          // Tail of double linked list for LFU
	LRUHead    *Node[KeyT, ValueT]          // Head of double linked list for LRU
	LRUTail    *Node[KeyT, ValueT]          // Tail of double linked list for LRU
	TTLHead    *Node[KeyT, ValueT]          // Head of double linked list for TTL
	TTLTail    *Node[KeyT, ValueT]          // Tail of double linked list for TTL
}

const ( // cache types
	Simple   = iota // No eviction
	LFUCache        // LFU eviction
	LRUCache        // LRU eviction
)

type Config struct { // cache parameters
	MaxSize uint64 // Maximum cache size, 0 - indefinite
	Type    uint   // One of Simple (no eviction on growth), LFUCache (least frequently used eviction), LRU (least recently used eviction)
	TTL     uint64 // Time to live for records in seconds, 0 - indefinite
}

func NewCache[KeyT comparable, ValueT any](ctx context.Context, config Config) (Cache[KeyT, ValueT], error) {
	return &cache[KeyT, ValueT]{
		Config:  config,
		Storage: make(map[KeyT]*Node[KeyT, ValueT]),
	}, nil
}

func (c *cache[KeyT, ValueT]) updateLFU(node *Node[KeyT, ValueT]) { // update node position in LFU double linked list towards LFUHead
	node.LFUCount++
	for node.LFUPrev != nil && node.LFUPrev.LFUCount < node.LFUCount {
		if node.LFUNext == nil {
			c.LFUTail = node.LFUPrev
		} else {
			node.LFUNext.LFUPrev = node.LFUPrev
		}
		node.LFUNext = node.LFUPrev
		node.LFUPrev = node.LFUPrev.LFUPrev
		if node.LFUPrev == nil {
			c.LFUHead = node
		}
		node.LFUNext.LFUPrev = node
	}
}

func (c *cache[KeyT, ValueT]) updateLRU(node *Node[KeyT, ValueT]) { // update node position in LRU double linked list, always to LRUHead
	if node.LRUPrev != nil {
		node.LRUPrev.LRUNext = node.LRUNext
		if node.LRUNext == nil {
			c.LRUTail = node.LRUPrev
		} else {
			node.LRUNext.LRUPrev = node.LRUPrev
		}
		node.LRUPrev = nil
		node.LRUNext = c.LRUHead
		c.LRUHead = node
	}
}

func (c *cache[KeyT, ValueT]) updateTTL(node *Node[KeyT, ValueT]) { // update node position in TTL double linked list, always to TTLHead
	node.UsedAt = time.Now()
	if node.LRUPrev != nil {
		node.LRUPrev.LRUNext = node.LRUNext
		if node.LRUNext == nil {
			c.LRUTail = node.LRUPrev
		} else {
			node.LRUNext.LRUPrev = node.LRUPrev
		}
		node.LRUPrev = nil
		node.LRUNext = c.LRUHead
		c.LRUHead = node
	}
}

func (c *cache[KeyT, ValueT]) insertLFU(node *Node[KeyT, ValueT]) { // insert new node to LFU double linked list, always at LFUTail
	node.LFUCount = 1
	node.LFUPrev = nil
	node.LFUNext = c.LFUHead
	c.LFUHead = node
}

func (c *cache[KeyT, ValueT]) insertLRU(node *Node[KeyT, ValueT]) { // insert new node to LFU double linked list, always at LRUHead
	node.LRUPrev = nil
	node.LRUNext = c.LRUHead
	c.LRUHead = node
}

func (c *cache[KeyT, ValueT]) insertTTL(node *Node[KeyT, ValueT]) { // insert new node to LFU double linked list, always at TTLHead
	node.UsedAt = time.Now()
	node.TTLPrev = nil
	node.TTLNext = c.LFUHead
	c.LFUHead = node
}

func (c *cache[KeyT, ValueT]) removeLFU(node *Node[KeyT, ValueT]) { // remove node from LFU double linked list
	node.LFUCount = 0
	if node.LFUNext == nil {
		c.LFUTail = node.LFUPrev
	} else {
		node.LFUNext.LFUPrev = node.LFUPrev
	}
	if node.LFUPrev == nil {
		c.LFUHead = node.LFUNext
	} else {
		node.LFUPrev.LFUNext = node.LFUNext
	}
	node.LFUPrev = nil
	node.LFUNext = nil
}

func (c *cache[KeyT, ValueT]) removeLRU(node *Node[KeyT, ValueT]) { // remove node from LRU double linked list
	if node.LRUNext == nil {
		c.LRUTail = node.LRUPrev
	} else {
		node.LRUNext.LRUPrev = node.LRUPrev
	}
	if node.LRUPrev == nil {
		c.LRUHead = node.LRUNext
	} else {
		node.LRUPrev.LRUNext = node.LRUNext
	}
	node.LRUPrev = nil
	node.LRUNext = nil
}

func (c *cache[KeyT, ValueT]) removeTTL(node *Node[KeyT, ValueT]) { // remove node from TTL double linked list
	node.UsedAt = time.Time{}
	if node.TTLNext == nil {
		c.TTLTail = node.TTLPrev
	} else {
		node.TTLNext.TTLPrev = node.TTLPrev
	}
	if node.TTLPrev == nil {
		c.TTLHead = node.TTLNext
	} else {
		node.TTLPrev.TTLNext = node.TTLNext
	}
	node.TTLPrev = nil
	node.TTLNext = nil
}

func (c *cache[KeyT, ValueT]) updateNode(node *Node[KeyT, ValueT], value ValueT) {
	c.Lock.Lock()
	c.Storage[node.Key].Value = value
	c.updateLFU(node)
	c.updateLRU(node)
	c.updateTTL(node)
	c.Lock.Unlock()
}

func (c *cache[KeyT, ValueT]) refreshNode(node *Node[KeyT, ValueT]) {
	c.Lock.Lock()
	c.updateLFU(node)
	c.updateLRU(node)
	c.Lock.Unlock()
}

func (c *cache[KeyT, ValueT]) insertNode(node *Node[KeyT, ValueT]) {
	c.Lock.Lock()
	c.Storage[node.Key] = node
	c.insertLFU(node)
	c.insertLRU(node)
	c.insertTTL(node)
	c.Size++
	c.Lock.Unlock()
}

func (c *cache[KeyT, ValueT]) removeNode(node *Node[KeyT, ValueT]) {
	delete(c.Storage, node.Key)
	c.removeLFU(node)
	c.removeLRU(node)
	c.removeTTL(node)
	c.Size--
}

func (c *cache[KeyT, ValueT]) evictByTTL() {
	c.Lock.Lock()
	for c.TTLTail != nil && uint64(time.Since(c.TTLTail.UsedAt).Seconds()) >= c.Config.TTL {
		c.removeNode(c.TTLTail)
	}
	c.Lock.Unlock()
}

func (c *cache[KeyT, ValueT]) evictByLFU() {
	c.Lock.Lock()
	for c.LFUTail != nil && c.Size >= c.Config.MaxSize {
		c.removeNode(c.LFUTail)
	}
	c.Lock.Unlock()
}

func (c *cache[KeyT, ValueT]) evictByLRU() {
	c.Lock.Lock()
	for c.LRUTail != nil && c.Size >= c.Config.MaxSize {
		c.removeNode(c.LRUTail)
	}
	c.Lock.Unlock()
}

func (c *cache[KeyT, ValueT]) Set(ctx context.Context, key KeyT, value ValueT) bool { // upsert value to cache, returns false if cache is full
	c.ConfigLock.RLock()
	defer c.ConfigLock.RUnlock()

	c.Lock.RLock()
	if node, ok := c.Storage[key]; ok { // update value in cache
		c.Lock.RUnlock()
		c.updateNode(node, value)
		return true
	}
	if c.Config.Type == Simple && c.Config.TTL > 0 { // evict only by TTL
		c.Lock.RUnlock()
		c.evictByTTL()
	} else if c.Config.MaxSize > 0 && c.Size >= c.Config.MaxSize { // need to evict some records to insert new one
		if c.Config.TTL > 0 {
			c.Lock.RUnlock()
			c.evictByTTL()
			c.Lock.RLock()
		}
		if c.Config.Type != Simple && c.Size >= c.Config.MaxSize {
			c.Lock.RUnlock()
			switch c.Config.Type {
			case LFUCache:
				c.evictByLFU()
			case LRUCache:
				c.evictByLRU()
			}
			c.Lock.RLock()
		}
	}
	if c.Config.MaxSize == 0 || c.Size < c.Config.MaxSize { // insert new record if not full
		node := Node[KeyT, ValueT]{
			Key:   key,
			Value: value,
		}
		c.Lock.RUnlock()
		c.insertNode(&node)
		return true
	}
	c.Lock.RUnlock()
	return false // cache is full and no records can be evicted
}

func (c *cache[KeyT, ValueT]) Get(ctx context.Context, key KeyT) (*ValueT, bool) { // return value from cache, if not in cache returns false
	c.ConfigLock.RLock()
	defer c.ConfigLock.RUnlock()

	c.Lock.RLock()
	if node, ok := c.Storage[key]; ok {
		if c.Config.TTL > 0 && uint64(time.Since(node.UsedAt).Seconds()) >= c.Config.TTL {
			c.Lock.RUnlock()
			c.evictByTTL()
			return nil, false
		}
		c.Lock.RUnlock()
		c.refreshNode(node)
		return &node.Value, true
	}
	c.Lock.RUnlock()
	return nil, false
}

func (c *cache[KeyT, ValueT]) Invalidate(ctx context.Context, key KeyT) bool { // removes record from cache, if not in cache returns false
	c.ConfigLock.RLock()
	defer c.ConfigLock.RUnlock()

	c.Lock.RLock()
	if node, ok := c.Storage[key]; ok {
		c.Lock.RUnlock()
		c.removeNode(node)
		return true
	}
	c.Lock.RUnlock()
	return false
}

func (c *cache[KeyT, ValueT]) SetConfig(ctx context.Context, config Config) error { // sets new configuration for cache at any time
	c.ConfigLock.Lock()
	defer c.ConfigLock.Unlock()

	c.Config = config
	return nil
}

func (c *cache[KeyT, ValueT]) Clear(ctx context.Context) error { // Clear cache, with hope that GC will do all dirty work for us
	c.ConfigLock.Lock()
	defer c.ConfigLock.Unlock()

	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Storage = make(map[KeyT]*Node[KeyT, ValueT])
	c.Size = 0
	c.LFUHead = nil
	c.LFUTail = nil
	c.LRUHead = nil
	c.LRUTail = nil
	c.TTLHead = nil
	c.TTLTail = nil
	return nil
}
