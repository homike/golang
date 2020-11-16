package channeltest

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

// 第一种方法其实有些问题, shutdown 如果不加锁, 没法完全阻止消息的写入
// 因此并不能在退出时, 清空掉channel中的buf
func Test_ChanClose1(t *testing.T) {
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

// 参考: https://zhuanlan.zhihu.com/p/32529039

// 当有多个消费者，多个生产者时，如何优雅的关闭go channel?
// 添加一个额外Channel, 用于关闭
func Test_ChannelClose2(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// ...
	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})
	// stopCh is an additional signal channel.
	// Its sender is the moderator goroutine shown below.
	// Its reveivers are all senders and receivers of dataCh.
	toStop := make(chan string, 1)
	// The channel toStop is used to notify the moderator
	// to close the additional signal channel (stopCh).
	// Its senders are any senders and receivers of dataCh.
	// Its reveiver is the moderator goroutine shown below.

	var stoppedBy string

	// moderator
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// senders
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(MaxRandomNumber)
				if value == 0 {
					// Here, a trick is used to notify the moderator
					// to close the additional signal channel.
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}

				// The first select here is to try to exit the goroutine
				// as early as possible. This select blocks with one
				// receive operation case and one default branches will
				// be optimized as a try-receive operation by the
				// official Go compiler.
				select {
				case <-stopCh:
					return
				default:
				}

				// Even if stopCh is closed, the first branch in the
				// second select may be still not selected for some
				// loops (and for ever in theory) if the send to
				// dataCh is also unblocked.
				// This is why the first select block is needed.
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()

			for {
				// Same as the sender goroutine, the first select here
				// is to try to exit the goroutine as early as possible.
				select {
				case <-stopCh:
					return
				default:
				}

				// Even if stopCh is closed, the first branch in the
				// second select may be still not selected for some
				// loops (and for ever in theory) if the receive from
				// dataCh is also unblocked.
				// This is why the first select block is needed.
				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == MaxRandomNumber-1 {
						// The same trick is used to notify
						// the moderator to close the
						// additional signal channel.
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	// ...
	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}
