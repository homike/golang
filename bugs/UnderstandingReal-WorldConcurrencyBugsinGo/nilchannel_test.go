package bugs

import (
	"fmt"
	"testing"
)

// 只能使用初始化后得channel, 发送数据给 nil channel 或者 从 nil channel接收数据
// 都会导致阻塞
func TestNilChannelSendOrRecv(t *testing.T) {
	var ch chan struct{}

	ch <- struct{}{}

	fmt.Println("step 1")

	<-ch

	fmt.Println("step 2")
}
