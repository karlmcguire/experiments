package simple

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

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

func TestSet(t *testing.T) {
	c := NewCache(8)
	fmt.Println(c.Set(1, 1, 5))
	c.Get(1)
	c.Get(1)
	c.Get(1)
	c.Get(1)
	c.Get(1)
	c.Get(1)
	c.Get(1)
	c.Get(1)
	fmt.Println(c.Set(2, 1, 1))
	c.Get(2)
	c.Get(2)
	c.Get(2)
	fmt.Println(c.Set(3, 1, 1))
	fmt.Println(c.Set(4, 1, 1))
	time.Sleep(time.Second)
	fmt.Println(c.Set(5, 1, 1))
}

func BenchmarkGet(b *testing.B) {
	s := NewCache(numKeys)
	for i := range keys {
		s.Set(keys[i], 0, 1)
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
	s := NewCache(numKeys)
	b.SetBytes(1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := rand.Int() & keyMask; pb.Next(); i++ {
			s.Set(keys[i&keyMask], 0, 1)
		}
	})
}
