package bugs

import (
	"fmt"
	"testing"
	"time"
)

func TestNoBlockTimer(t *testing.T) {
	dur := 10 * time.Second
	timer := time.NewTimer(0)
	if dur > 0 {
		timer = time.NewTimer(dur)
	}
	select {
	case <-timer.C:
	}

	fmt.Println("timer")
}
