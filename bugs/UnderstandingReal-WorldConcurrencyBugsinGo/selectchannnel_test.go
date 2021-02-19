package bugs

import (
	"testing"
	"time"
)

// 当同时监听两个channel时候, 无法保证接收到channel消息的时序
// 因此只能在之前写读取stop channel
func TestNoBlockSelectTwoChannel(t *testing.T) {
	chStop := make(chan struct{})
	ticker := time.NewTicker(time.Second)

	for {
		// 正确写法
		//select {
		//case <-chStop:
		//	return
		//default:
		//}

		select {
		case <-chStop:
			return
		case <-ticker.C:
		}
	}
}
