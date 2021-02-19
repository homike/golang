package bugs

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

// 当 阻塞channel  和 mutex一起使用时, 容易导致死锁, 需要非常注意
// 当向 阻塞channel中写入数据时, 也需要非常注意，
// 要么使用select，要么换成非阻塞channel
func TestBlockByChannelAndMutex(t *testing.T) {
	var m sync.Mutex
	ch := make(chan struct{})

	fmt.Println("Start Gorutines: ", runtime.NumGoroutine())

	go func() {
		m.Lock()
		defer m.Unlock()

		<-ch
		// 正确写法
		//select {
		//case <-ch:
		//default:
		//}

	}()

	go func() {
		for {
			m.Lock()
			defer m.Unlock()
			ch <- struct{}{}
		}
	}()

	fmt.Println("End Gorutines: ", runtime.NumGoroutine())
}
