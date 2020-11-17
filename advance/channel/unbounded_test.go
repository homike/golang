package channel

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func Test_Unbounded(t *testing.T) {
	bound := NewUnbounded()

	go func() {
		for i := 0; i < math.MaxInt32; i++ {
			bound.Put(i)
			bound.Put(i)
			bound.Put(i)
			time.Sleep(1 * time.Second)
		}
	}()

	for i := 0; i < 2; i++ {
		go func(index int) {
			for {
				select {
				case t := <-bound.Get():
					bound.Load()
					v := t.(int)
					fmt.Println("index: ", index, ", len: ", len(bound.backlog), "value: ", v)
					time.Sleep(1 * time.Second)
				}
			}
		}(i)
	}

	select {
	case <-time.After(3 * time.Second):
	}
}
