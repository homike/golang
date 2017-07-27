package locktest

import (
	"sync"
	"time"
)

type LockTest struct {
	numOfClient uint

	WGroup sync.WaitGroup
	lock   sync.RWMutex
}

func (lockTest *LockTest) StartTest() {
	lockTest.numOfClient = 10

	for i := 0; i < 10; i++ {
		go func(idx int) {
			lockTest.WGroup.Add(1)
			for {
				time.Sleep(2)
			}
			lockTest.WGroup.Done()
		}(i)
	}
}

func (lockTest *LockTest) Reload() {

	lockTest.lock.Lock()
	defer lockTest.lock.Unlock()
}
