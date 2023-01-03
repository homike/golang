package base

import (
	"fmt"
	"testing"
	"unsafe"
)

type T1 struct {
	a int
	b int
	c int
}

type T2 struct {
	*T1
}

type T3 struct {
	t *T1
}

func TestSizeof(t *testing.T) {
	t1 := &T1{a: 1, b: 1, c: 1}
	t2 := &T2{T1: t1}
	t3 := &T3{t: t1}
	fmt.Println("size1: ", unsafe.Sizeof(t2))
	fmt.Println("size2: ", unsafe.Sizeof(t3))
}
