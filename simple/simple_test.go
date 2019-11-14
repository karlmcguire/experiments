package simple

import (
	"fmt"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	c := NewCache(4)
	fmt.Println(c.Set(1, 1, 1))
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
