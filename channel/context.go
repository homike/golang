package channeltest

import (
	"context"
	"fmt"
	"time"
)

// context.Background(): 所有context的root, 不能被cancel
// context.WithCancel(): 返回一个继承的Context, 这个Context在父Context的Done被关闭时关闭自己的Done通道，或者在自己被Cancel的时候关闭自己的Done。
// 					 	 WithCancel同时还返回一个取消函数cancel，这个cancel用于取消当前的Context。
func childFunc(cont context.Context, num *int) {
	ctx, _ := context.WithCancel(cont)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("child Done:", ctx.Err())
			return
		}
	}
}

func RunContext() {
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("parent one:", ctx.Err())
					return
				case dst <- n:
					n++
					go childFunc(ctx, &n)
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	for n := range gen(ctx) {
		fmt.Println(n)
		if n >= 5 {
			break
		}
	}
	cancel()
	time.Sleep(5 * time.Second)
}
