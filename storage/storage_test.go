package storage

import (
	"math/rand"
	"testing"

	"github.com/dgraph-io/ristretto/sim"
)

const (
	numKeys = 1e7
	keyMask = numKeys - 1
)

var (
	keys [numKeys]uint64
)

func init() {
	zipf := sim.NewZipfian(1.25, 2, numKeys)
	for i := range keys {
		keys[i], _ = zipf()
	}
}

func BenchmarkGet(b *testing.B) {
	s := NewUnified(numKeys)
	for i := range keys {
		s.Set(keys[i], 0, 1, int64(rand.Int()))
	}
	b.SetBytes(1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := rand.Int() & keyMask; pb.Next(); i++ {
			s.Get(keys[i&keyMask])
		}
	})
}

func BenchmarkSet(b *testing.B) {
	s := NewUnified(numKeys)
	b.SetBytes(1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := rand.Int() & keyMask; pb.Next(); i++ {
			s.Set(keys[i&keyMask], 0, 1, 0)
		}
	})
}
