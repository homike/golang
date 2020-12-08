package gobugs

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 1.gorutine的参数传递, 同步参数列表传递，不要直接使用局部变量，会导致意想不到的问题
func gorutineArgs() {
	for i := 17; i <= 21; i++ {
		/*
			go func() {
				_ := fmt.Sprintf("%d", i)
			}()
		*/
		fmt.Println("start1")
		go func(v int) {
			fmt.Println("end")
			fmt.Printf("%d \n", v)
		}(i)
	}
}

// 1.当stopCh 和 ticker 同时触发时，可能不能正常的return, 所以需要在循环开头，加上判断
func f() {}
func StopCh(stopCh chan int) {
	ticker := time.NewTicker(10)
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		f()

		select {
		case <-stopCh:
			return
		case <-ticker.C:
		default:
		}
	}

}

// 1.select 等待多个channnel时，都准备就绪的情况，会随机选择一个。
// 需求: 如果dur 为 0, 则没有超时, 一直等待done
//		 如果等待timer.C, 则dur为0的时候，就直接返回了
func SelectTimeOut(ctx context.Context, dur int) {
	//timer := time.NewTimer(time.Duration(0))
	var timeout <-chan time.Time
	if dur > 0 {
		//timer = time.NewTimer(time.Duration(dur))
		timeout = time.NewTimer(time.Duration(dur)).C
	}

	select {
	//case <-timer.C:
	case <-timeout:
		fmt.Println("timeount")
		return
	case <-ctx.Done():
		fmt.Println("done")
		return
	}
}

func TestSelectTimeOut(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	SelectTimeOut(ctx, 0)
}
