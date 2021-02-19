package bugs

import (
	"fmt"
	"testing"
)

// 关闭已经关闭的channel导致的panic
func TestNoBlockByCloseTwice(t *testing.T) {
	ch := make(chan struct{})

	for i := 0; i < 3; i++ {
		go func() {
			select {
			case <-ch:
				fmt.Println("ch")
			default:
				fmt.Println("close")
				close(ch)
				// 正确做法
				//sync.Once.Do(func(){
				//	close(ch)
				//})
			}
		}()
	}

	for {
	}
}
