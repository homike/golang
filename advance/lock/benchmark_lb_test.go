package lock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type locklb struct {
	arr  []int
	next int
	sync.RWMutex
}

/*
func Benchmark_roundlb(b *testing.B) {
	l := &locklb{
		arr:  make([]int, 100, 100),
		next: 0,
	}
	b.SetParallelism(10000)

	// 重置计时器
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Lock()

			l.arr[l.next]++
			l.next = (l.next + 1) % len(l.arr)

			l.Unlock()
		}
	})
}

func Benchmark_randlb(b *testing.B) {
	l := &locklb{
		arr:  make([]int, 100, 100),
		next: 0,
	}
	b.SetParallelism(10000)

	// 重置计时器
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.RLock()

			rIndex := rand.Intn(len(l.arr))
			_ = l.arr[rIndex]

			l.RUnlock()
		}
	})
}
*/

func Test_roundlb(t *testing.T) {
	l := &locklb{
		arr:  make([]int, 100, 100),
		next: 0,
	}

	start := make(chan struct{})
	end := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			select {
			case <-start:
				for {
					select {
					case <-end:
						wg.Done()
						return
					default:
					}

					l.Lock()
					l.arr[l.next]++
					l.next = (l.next + 1) % len(l.arr)
					l.Unlock()
				}
			}
		}()
	}
	close(start)

	time.Sleep(10 * time.Second)

	close(end)
	wg.Wait()

	min, max := l.arr[0], l.arr[0]
	for k, v := range l.arr {
		fmt.Println("K: ", k, " v: ", v)
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	fmt.Println("diff : ", max-min)

}
