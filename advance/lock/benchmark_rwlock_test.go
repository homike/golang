package lock

import (
	"sync"
	"testing"
)

type lockrw struct {
	next int
	sync.RWMutex
	mu sync.Mutex
}

func Benchmark_nolock(b *testing.B) {
	l := &lockrw{
		next: 0,
	}
	b.SetParallelism(10000)

	// 重置计时器
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.next++
		}
	})
}

func Benchmark_rlock(b *testing.B) {
	l := &lockrw{
		next: 0,
	}
	b.SetParallelism(10000)

	// 重置计时器
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.RLock()

			l.next++

			l.RUnlock()
		}
	})
}

func Benchmark_rwlock(b *testing.B) {
	l := &lockrw{
		next: 0,
	}
	b.SetParallelism(10000)

	// 重置计时器
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Lock()

			l.next++

			l.Unlock()
		}
	})
}

func Benchmark_mutexlock(b *testing.B) {
	l := &lockrw{
		next: 0,
	}
	b.SetParallelism(10000)

	// 重置计时器
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.mu.Lock()

			l.next++

			l.mu.Unlock()
		}
	})
}
