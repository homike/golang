package locktest

import (
	"fmt"
	"sync"
	"time"
)

func RunTest() {
	var num = 50000000
	var c = make(chan int, 3000)

	var rwMutex = newRwmutex()
	var w = &sync.WaitGroup{}
	w.Add(num)
	t1 := time.Now()
	for i := 0; i < num; i++ {
		c <- 0
		go func(index int) {
			defer w.Done()
			_ = rwMutex.get(index)
			<-c
		}(i)
	}
	w.Wait()
	t2 := time.Now()

	var mutex = newMutex()
	w.Add(num)
	t3 := time.Now()
	for i := 0; i < num; i++ {
		c <- 0
		go func(index int) {
			defer w.Done()
			_ = mutex.get(index)
			<-c
		}(i)
	}
	w.Wait()
	t4 := time.Now()
	fmt.Println("rwmutex cost:", t2.Sub(t1).String())
	fmt.Println("mutex cost:", t4.Sub(t3).String())
}

type rwmutex struct {
	mu    *sync.RWMutex
	ipmap map[int]int
}

type mutex struct {
	mu    *sync.Mutex
	ipmap map[int]int
}

func (t *rwmutex) get(i int) int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.ipmap[i]
}

func (t *mutex) get(i int) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.ipmap[i]
}

func newRwmutex() *rwmutex {
	var t = &rwmutex{}
	t.mu = &sync.RWMutex{}
	t.ipmap = make(map[int]int, 100)

	for i := 0; i < 100; i++ {
		t.ipmap[i] = 0
	}
	return t
}

func newMutex() *mutex {
	var t = &mutex{}
	t.mu = &sync.Mutex{}
	t.ipmap = make(map[int]int, 100)

	for i := 0; i < 100; i++ {
		t.ipmap[i] = 0
	}
	return t
}
