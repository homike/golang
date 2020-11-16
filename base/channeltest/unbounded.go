package channeltest

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type Unbounded struct {
	c       chan interface{}
	mu      sync.Mutex
	backlog []interface{}
}

func NewUnbounded() *Unbounded {
	return &Unbounded{c: make(chan interface{})}
}

func (b *Unbounded) Put(t interface{}) {
	b.mu.Lock()
	if len(b.backlog) == 0 {
		select {
		case b.c <- t:
			b.mu.Unlock()
			fmt.Println("put 1")
			return
		default:
			fmt.Println("put 2")
		}
	}
	fmt.Println("put 3")
	b.backlog = append(b.backlog, t)
	b.mu.Unlock()
}

func (b *Unbounded) Load() {
	b.mu.Lock()
	if len(b.backlog) > 0 {
		select {
		case b.c <- b.backlog[0]:
			b.backlog[0] = nil
			b.backlog = b.backlog[1:]
			fmt.Println("load 1")
		default:
			fmt.Println("load 2")
		}
	}
	b.mu.Unlock()
}

func (b *Unbounded) Get() <-chan interface{} {
	fmt.Println("Get")
	return b.c
}

//------------------------------------------------------
func RunUnbounded() {
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

	for {
	}
}
