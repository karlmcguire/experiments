package cache

import (
	"sync"

	"github.com/karlmcguire/plrum"
)

type item struct {
	key uint64
	val []byte
}

type Cache struct {
	sync.Mutex
	keys map[uint64]uint64
	data []item
	meta plrum.Policy
	mask uint64
	used uint64
}

func NewCache(size uint64) *Cache {
	return &Cache{
		keys: make(map[uint64]uint64, size),
		data: make([]item, size),
		meta: plrum.NewPolicy(size),
		mask: size - 1,
	}
}

func (c *Cache) Get(key uint64) []byte {
	c.Lock()
	defer c.Unlock()
	block := c.keys[key]
	c.meta.Hit(block)
	return c.data[block].val
}

func (c *Cache) Set(key uint64, val []byte) (victim uint64) {
	c.Lock()
	defer c.Unlock()
	// if already exists, just update
	if block, ok := c.keys[key]; ok {
		c.data[block].val = val
		c.meta.Hit(block)
		return
	}
	// find a new open block
	block := c.meta.Evict()
	if c.used > c.mask {
		victim = c.data[block].key
		delete(c.keys, victim)
		c.used--
	}
	c.keys[key] = block
	c.data[block] = item{key, val}
	c.meta.Hit(block)
	c.used++
	return
}
