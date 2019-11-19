package cache

import (
	"fmt"
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
	block, ok := c.keys[key]
	if !ok {
		return nil
	}
	c.meta.Hit(block)
	return c.data[block].val
}

func (c *Cache) Set(key uint64, val []byte) (victim uint64) {
	c.Lock()
	defer c.Unlock()
	// if already exists, just update
	if block, ok := c.keys[key]; ok {
		c.meta.Hit(block)
		c.data[block].val = val
		return
	}
	// find a new open block
	block := c.meta.Evict()
	c.meta.Hit(block)
	if c.used > c.mask {
		victim = c.data[block].key
		delete(c.keys, victim)
		c.used--
	}
	c.keys[key] = block
	c.data[block] = item{key, val}
	c.used++
	return
}

func (c *Cache) String() string {
	var out string
	out = "[ "
	for i := uint64(0); i <= c.mask; i++ {
		if i%8 == 0 {
			out += "\n"
		}
		out += fmt.Sprintf("%2d: %d  ", i, c.data[i].key)
	}
	return out + "]"
}
