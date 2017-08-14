package channeltest

import (
	"context"
	"fmt"
	"time"
)

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
