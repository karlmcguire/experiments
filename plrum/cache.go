package cache

import (
	"sync"

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

func (c *Cache) Get(key uint64) []byte {
	id := key & c.mask
	c.RLock()
	val := c.data[id]
	c.meta.Hit(id)
	c.RUnlock()
	return val
}

func (c *Cache) Set(key uint64, val []byte) uint64 {
	id := key & c.mask
	c.Lock()
	victim := c.meta.Evict()
	c.data[victim] = nil
	c.data[id] = val
	c.Unlock()
	return victim
}
