package channeltest

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// go1.7之后加入, 为了方便gorutine的管理
/*
	Context,为了多个gorutine之间共享状态和数据.
*/

/*
	context.Background(): 所有context的root, 不能被cancel

	context.WithCancel(): 返回一个继承的Context, 这个Context在父Context的Done被关闭时关闭自己的Done通道，或者在自己被Cancel的时候关闭自己的Done。
 					 	  WithCancel同时还返回一个取消函数cancel，这个cancel用于取消当前的Context。

	WithTimeout func(parent Context, timeout time.Duration) (Context, CancelFunc): 设置ctx的超时
	WithTimeout 等价于 WithDeadline(parent, time.Now().Add(timeout)).

	"context deadline exceeded"就是ctx超时的时候ctx.Err的错误消息。

	WithValue()函数将key和value保存在新的上下文对象中并返回该对象
*/
func childFunc(cont context.Context, num *int) {
	ctx, _ := context.WithCancel(cont)
	v, ok := ctx.Value("value").(int)
	fmt.Println("value", v, ",ok", ok)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("child Done:", ctx.Err())
			return
		}
	}
}

func Test_Context(t *testing.T) {
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

	ctx1, cancel := context.WithCancel(context.Background())
	//ctx1, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	//ctx1, cancel := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))

	ctx := context.WithValue(ctx1, "value", 123456)
	for n := range gen(ctx) {
		fmt.Println(n)
		if n >= 5 {
			break
		}
	}
	cancel()
	time.Sleep(5 * time.Second)
}
