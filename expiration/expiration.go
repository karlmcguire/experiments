package expiration

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	times    *list.List
	data     map[uint64][2]uint64
	used     uint64
	capacity uint64
}

func NewCache(capacity uint64) *Cache {
	return &Cache{
		data:     make(map[uint64][2]uint64),
		times:    list.New(),
		capacity: capacity,
	}
}

func (c *Cache) Get(key uint64) uint64 {
	c.RLock()
	val := c.data[key]
	c.RUnlock()
	return val[0]
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
	// check if items with set TTLs have expired first
	for e := c.times.Back(); e != nil; e = e.Prev() {
		item := e.Value.([3]uint64)
		if item[2] > uint64(time.Now().Unix()) {
			break
		}
		c.times.Remove(e)
		delete(c.data, item[0])
		c.used -= item[1]
		victims = append(victims, item[0])
	}
	for ; room < cost; room = c.capacity - c.used {
		// regular cost-based eviction
		/*
			i := 0
			minKey, minCost := uint64(0), uint64(0)
			for k, v := range c.data {
				// TODO
				if i++; i == 5 {
					break
				}
			}
		*/
	}
add:
	c.data[key] = [2]uint64{val, cost}
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
