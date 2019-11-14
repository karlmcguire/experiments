package simple

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	data     *sync.Map
	used     uint64
	capacity uint64
}

type item struct {
	hits    uint64
	cost    uint64
	created int64
	val     interface{}
}

func (i *item) priority() float64 {
	return float64(i.hits) / float64(time.Now().Unix()-i.created)
}

func NewCache(capacity uint64) *Cache {
	return &Cache{
		data:     &sync.Map{},
		capacity: capacity,
	}
}

func (c *Cache) Get(key uint64) interface{} {
	val, ok := c.data.Load(key)
	if !ok || val == nil {
		return nil
	}
	data := val.(*item)
	data.hits++
	return data.val
}

func (c *Cache) Set(key uint64, val interface{}, cost uint64) []uint64 {
	// check if new item is larger than capacity of the entire cache
	if cost > c.capacity {
		return nil
	}
	room := uint64(0)
	victims := make([]uint64, 0)
	// check if there's room for the new item without any evictions
	if room = c.capacity - atomic.LoadUint64(&c.used); room > cost {
		goto add
	}
	// continuously sample and evict items until there's enough room
	for ; room < cost; room = c.capacity - atomic.LoadUint64(&c.used) {
		i := 0
		minKey, minItem := uint64(0), &item{}
		c.data.Range(func(k, v interface{}) bool {
			candidate := v.(*item)
			fmt.Printf("%d : %.8f\n", k.(uint64), candidate.priority())
			if i == 0 || candidate.priority() < minItem.priority() {
				minKey, minItem = k.(uint64), candidate
			}
			if i++; i == 5 {
				return false
			}
			return true
		})
		c.data.Delete(minKey)
		atomic.AddUint64(&c.used, ^uint64(minItem.cost-1))
		victims = append(victims, minKey)
	}
add:
	c.data.Store(key, &item{1, cost, time.Now().Unix(), val})
	atomic.AddUint64(&c.used, cost)
	return victims
}
