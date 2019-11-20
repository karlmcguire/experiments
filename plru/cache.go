package cache

import (
	"sync"

	"github.com/karlmcguire/plru"
)

const (
	numShards = 256
	shardMask = numShards - 1
)

type item struct {
	key uint64
	val []byte
}

type Cache [numShards]*shard

func NewCache(size uint64) Cache {
	var c Cache
	for i := range c {
		c[i] = newShard((size + shardMask) / numShards)
	}
	return c
}

func (c Cache) Get(key uint64) []byte {
	return c[key&shardMask].get(key)
}

func (c Cache) Set(key uint64, val []byte) {
	c[key&shardMask].set(key, val)
}

type shard struct {
	sync.RWMutex
	meta *plru.Policy
	keys map[uint64]uint64
	data []item
	used uint64
	mask uint64
}

func newShard(size uint64) *shard {
	return &shard{
		meta: plru.NewPolicy(size),
		keys: make(map[uint64]uint64, size),
		data: make([]item, size),
		mask: size - 1,
	}
}

func (s *shard) get(key uint64) []byte {
	s.RLock()
	key, ok := s.keys[key]
	s.RUnlock()
	if !ok {
		return nil
	}
	s.meta.Hit(key)
	return s.data[key].val
}

func (s *shard) set(key uint64, val []byte) (victim uint64) {
	s.Lock()
	defer s.Unlock()
	// if already exists, just update
	if block, ok := s.keys[key]; ok {
		s.meta.Hit(block)
		s.data[block].val = val
		return
	}
	// find a new open block
	block := s.meta.Evict()
	// check if eviction is needed
	if s.used > s.mask || s.data[block].val != nil {
		victim = s.data[block].key
		delete(s.keys, victim)
		goto add
	}
	s.used++
add:
	s.keys[key] = block
	s.data[block] = item{key, val}
	s.meta.Hit(block)
	return
}
