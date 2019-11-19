package cache

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	c := NewCache(64)

	for i := uint64(0); i < 64; i++ {
		c.Set(i, []byte("data"))
	}

	fmt.Printf("used: %d\n", len(c.keys))
	fmt.Printf("%064b\n", c.meta[0])

	fmt.Println(c.Set(3255, []byte("data")))
	fmt.Printf("used: %d\n", len(c.keys))
	fmt.Printf("%064b\n", c.meta[0])

	fmt.Println(c.Set(32, []byte("data")))
	fmt.Printf("used: %d\n", len(c.keys))
	fmt.Printf("%064b\n", c.meta[0])

}
