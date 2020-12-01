package main

import (
	"fmt"
	"testing"
)

type worker interface {
	work()
}

type person struct {
	name string
	worker
}

// 这个能测试通过吗?
func TestQuestion1(t *testing.T) {
	var w worker = person{}
	fmt.Println(w)
}
