package sync

import (
	"fmt"
	"testing"

	"go.uber.org/atomic"
)

func TestAtomicCAS(t *testing.T) {
	var b atomic.Bool
	fmt.Println("init 0 ", b.Load())

	fmt.Println("ret 1 ", b.CAS(false, true))
	fmt.Println("ret 1 ", b.CAS(false, true))
	fmt.Println("ret 1 ", b.CAS(false, true))
	fmt.Println("ret 1 ", b.CAS(false, true))
	fmt.Println("ret 1 ", b.CAS(false, true))

	fmt.Println("---------------")
	b.Store(false)
	fmt.Println("ret 1 ", b.CAS(false, true))
	fmt.Println("ret 1 ", b.CAS(false, true))
	fmt.Println("ret 1 ", b.CAS(false, true))
	fmt.Println("ret 1 ", b.CAS(false, true))
	fmt.Println("init 0 ", b.Load())

	//fmt.Println("ret 2 ", b.CAS(true, false))
	//fmt.Println("ret 2 ", b.CAS(false, false))
}
