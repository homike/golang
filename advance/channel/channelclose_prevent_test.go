package channel

import (
	"fmt"
	"testing"
)

func Test_ClosePreven(t *testing.T) {
	ch := make(chan struct{})

	close(ch)

	select {
	case <-ch:
		fmt.Println("close 111")
	default:
		fmt.Println("close 222")
		close(ch)
	}
}
