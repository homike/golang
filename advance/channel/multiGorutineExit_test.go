package channel

import (
	"sync"
	"testing"
)

//---------------------------------------------------------
// 测试以下两个gorutine退出方式的差别
//	1. 全局channel 通知所有gorutine退出
//  2. 通过gorutine私有的channel控制退出
//---------------------------------------------------------
var chanExit chan struct{}
var WGroup sync.WaitGroup

type actor struct {
	chanExit chan struct{}
}

func NewActor() *actor {
	return &actor{
		chanExit: make(chan struct{}),
	}
}

func (a *actor) Start() {
	WGroup.Add(1)

	go func() {
		defer WGroup.Done()

		for {
			select {
			case <-chanExit:
				return
			case <-a.chanExit:
				return
			}
		}
	}()
}

func (a *actor) Stop() {
	close(a.chanExit)
}

func Benchmark_MultiGorutineExit_GlobalChannel(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()

	chanExit = make(chan struct{})
	var actors sync.Map

	// 开始计时器
	b.StartTimer()

	for i := 0; i < 10000; i++ {
		a := NewActor()
		actors.LoadOrStore(i, a)
		a.Start()
	}

	close(chanExit)
	WGroup.Wait()
}

func Benchmark_MultiGorutineExit_PerChannel(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()

	chanExit = make(chan struct{})
	var actors sync.Map

	// 开始计时器
	b.StartTimer()

	for i := 0; i < 10000; i++ {
		a := NewActor()
		actors.LoadOrStore(i, a)
		a.Start()
	}

	//fmt.Println("============exit")
	actors.Range(
		func(_, v interface{}) bool {
			a := v.(*actor)
			a.Stop()
			return true
		})
	WGroup.Wait()
}
