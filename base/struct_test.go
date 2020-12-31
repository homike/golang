package base

import (
	"fmt"
	"testing"
)

type emptyStruct struct{}

// 空结构体的变量的内存地址都是一样的
func TestStruct(t *testing.T) {
	a := struct{}{}
	b := struct{}{}
	c := emptyStruct{}

	fmt.Printf("a: %p, b: %p, c: %p", &a, &b, &c)
}
