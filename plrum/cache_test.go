package cache

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	c := NewCache(128)
	rounds := uint64(128)
	for i := uint64(0); i < rounds; i++ {
		c.Set(i, []byte(fmt.Sprintf("%d", i)))
	}
	fmt.Println(c)
	for i := uint64(0); i < rounds; i++ {
		val := c.Get(i)
		if val == nil {
			continue
		}
		if string(val) != fmt.Sprintf("%d", i) {
			fmt.Printf("get %d: %s\n", i, string(c.Get(i)))
			//t.Fatal("key != value")
		}
	}
}

func BenchmarkGet(b *testing.B) {
	c := NewCache(64)
	c.Set(0, []byte("data"))
	b.SetBytes(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(0)
		}
	})
}

func BenchmarkSet(b *testing.B) {
	c := NewCache(64)
	b.SetBytes(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(0, []byte("data"))
		}
	})
}
