package cache

import (
	"sync"

	"github.com/cespare/xxhash"
	"github.com/karlmcguire/plrum"
)

type Cache struct {
	sync.RWMutex
	data [][]byte
	meta plrum.Policy
	mask uint64
	used uint64
}

func NewCache(size uint64) *Cache {
	return &Cache{
		data: make([][]byte, size),
		meta: plrum.NewPolicy(size),
		mask: size - 1,
	}
}

func (c *Cache) Get(key []byte) []byte {
	id := xxhash.Sum64(key) & c.mask
	c.RLock()
	defer c.RUnlock()
	val := c.data[id]
	c.meta.Hit(id)
	return val
}

func (c *Cache) Set(key []byte, val []byte) uint64 {
	id := xxhash.Sum64(key) & c.mask
	c.Lock()
	defer c.Unlock()
	if c.used > c.mask {
		victim := c.meta.Evict()
		c.data[victim] = nil
		id = victim
		c.used--
	}
	c.data[id] = val
	c.meta.Hit(id)
	c.used++
	return 0
}
