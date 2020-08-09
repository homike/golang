package stack

import "testing"

func TestStack(t *testing.T) {
	PrintGoroutineMemConsume()
	/*
		for i := 0; i < 4; i++ {
			go func(index int) {
				t.Logf("%d, val: %s", index, PrintStack())
			}(i)
		}

		time.Sleep(1 * time.Second)
	*/
}
