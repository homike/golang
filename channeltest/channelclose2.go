package channeltest

import (
	"fmt"
	"time"
)

func ChanClose2() {

	msglist := make(chan int, 10)
	shutdown := false
	closeChan := make(chan struct{})

	// 模拟上行消息
	go func() {
		i := 0
		for {
			i++
			if shutdown {
				return
			}
			select {
			case msglist <- i:
			default:
				//time.Sleep(100 * time.Millisecond)
			}

		}
	}()

	// 模拟接收消息
	go func() {
	loop:
		for {
			select {
			case msg := <-msglist:
				_ = msg
				//fmt.Println("recv msg ", msg)
			case <-closeChan:
				shutdown = true
				break loop
			}
		}

		// 退出, 清空消息队列
		// close(msglist), 不要再接收端close这个队列
		for {
			select {
			case msg := <-msglist:
				fmt.Println("close msg ", msg)
			default:
				return
			}
		}
	}()

	time.Sleep(1 * time.Second)

	close(closeChan)

	time.Sleep(1 * time.Second)
	fmt.Println("exit")
}
