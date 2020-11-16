package lock

import (
	"sync"
	"time"
)

var (
	w    sync.WaitGroup
	lock sync.RWMutex
)

func Test_WaitGroup() {
	for i := 0; i < 10; i++ {
		go func(idx int) {
			w.Add(1)
			defer w.Done()
			for {
				time.Sleep(2)
			}
		}(i)
	}
}

func Reload() {
	lock.Lock()
	defer lock.Unlock()
}

var (
	a1   string
	once sync.Once
)

func setup() {
	a1 = "hello"
}

func doprint() {
	once.Do(setup)
	print(a1)
}

func Test_SyncOnce() {
	go doprint()

	go doprint()
}
