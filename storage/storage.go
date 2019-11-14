package storage

import "sync"

type Unified struct {
	shards [numShards]*shard
}

func NewUnified(size uint64) *Unified {
	s := &Unified{}
	for i := range s.shards {
		s.shards[i] = newShard(size / numShards)
	}
	return s
}

func (s *Unified) Get(key uint64) interface{} {
	val, _, _ := s.shards[key&shardMask].Get(key)
	return val
}

func (s *Unified) Set(key uint64, val interface{}, cost, ttl int64) []uint64 {
	s.shards[key&shardMask].Set(key, val, cost, ttl)
	return nil
}

func (s *Unified) Del(key uint64) {
	s.shards[key&shardMask].Del(key)
}

const (
	numShards = 256
	shardMask = numShards - 1
)

type shard struct {
	sync.RWMutex
	data map[uint64]shardItem
}

type shardItem struct {
	val  interface{}
	cost int64
	ttl  int64
	hits uint64
}

func newShard(size uint64) *shard {
	return &shard{
		data: make(map[uint64]shardItem, size),
	}
}

func (s *shard) Get(key uint64) (interface{}, int64, int64) {
	s.RLock()
	item := s.data[key]
	s.RUnlock()
	return item.val, item.cost, item.ttl
}

func (s *shard) Set(key uint64, val interface{}, cost, ttl int64) {
	s.Lock()
	s.data[key] = shardItem{val, cost, ttl}
	s.Unlock()
}

func (s *shard) Del(key uint64) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}
