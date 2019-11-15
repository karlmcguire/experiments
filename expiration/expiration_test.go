package expiration

import (
	"fmt"
	"testing"
	"time"
)

func TestExpiration(t *testing.T) {
	c := NewCache(10)
	c.Set(1, 1, 5, 1)
	c.Set(2, 2, 4, 1)
	time.Sleep(time.Second * 3)

	fmt.Println(c.Set(3, 3, 3, 0))
}
