package cache

import (
	"fmt"
	"sync"

	"github.com/karlmcguire/plru"
)

type item struct {
	key uint64
	val []byte
}

type Cache struct {
	sync.RWMutex
	keys map[uint64]uint64
	data []item
	meta *plru.Policy
	mask uint64
	used uint64
}

func NewCache(size uint64) *Cache {
	return &Cache{
		keys: make(map[uint64]uint64, size+1),
		data: make([]item, size+1),
		meta: plru.NewPolicy(size),
		mask: size - 1,
	}
}

func (c *Cache) Get(key uint64) []byte {
	c.RLock()
	block, ok := c.keys[key]
	c.RUnlock()
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
	fmt.Println(len(c.data), block)
	// check if eviction is needed
	if c.used > c.mask || c.data[block].val != nil {
		victim = c.data[block].key
		fmt.Printf("evicting %d\n", c.data[block].key)
		delete(c.keys, victim)
		goto add
	}
	c.used++
add:
	c.keys[key] = block
	c.data[block] = item{key, val}
	c.meta.Hit(block)
	return
}

func (c *Cache) String() string {
	var out string
	out = "[ "
	for i := uint64(0); i <= c.mask; i++ {
		if i%8 == 0 {
			out += "\n"
		}
		out += fmt.Sprintf("%3d: %3d  ", i, c.data[i].key)
	}
	return out + "\n]"
}
