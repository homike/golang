package channel

import (
	"fmt"
	"sync"
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
