package expiration

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	times    *list.List
	data     map[uint64]uint64
	used     uint64
	capacity uint64
}

func NewCache(capacity uint64) *Cache {
	return &Cache{
		data:     make(map[uint64]uint64),
		times:    list.New(),
		capacity: capacity,
	}
}

func (c *Cache) Get(key uint64) uint64 {
	c.RLock()
	val := c.data[key]
	c.RUnlock()
	return val
}

func (c *Cache) Set(key, val, cost, ttl uint64) []uint64 {
	c.Lock()
	defer c.Unlock()
	if cost > c.capacity {
		return nil
	}
	room := uint64(0)
	victims := make([]uint64, 0)
	if room = c.capacity - c.used; room > cost {
		goto add
	}
	for ; room < cost; room = c.capacity - c.used {
		// check if items with set TTLs have expired first
		for e := c.times.Back(); e != nil; e = e.Prev() {
			item := e.Value.([3]uint64)
			if item[2] < uint64(time.Now().Unix()) {
				c.times.Remove(e)
				delete(c.data, item[0])
				c.used -= item[1]
				victims = append(victims, item[0])
			}
		}
		// regular cost-based eviction
		// TODO
	}
add:
	c.data[key] = val
	c.used += cost
	if ttl > 0 {
		c.times.PushFront([3]uint64{
			key,
			cost,
			uint64(time.Now().Unix() + int64(ttl)),
		})
	}
	return victims
}
