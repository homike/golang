package main

import (
	"testing"
)

// channel 的哪些操作会引发 panic？

// 1. 关闭一个 nil 值 channel 会引发 panic
func TestCloseNilChannel(t *testing.T) {
	var ch chan struct{}
	close(ch)
}

// 2. 关闭一个已关闭的 channel 会引发 panic
func TestCloseClosedChannel(t *testing.T) {
	ch := make(chan struct{})
	close(ch)
	close(ch)
}

// 3. 向一个已关闭的 channel 发送数据
func TestSendClosedChannel(t *testing.T) {
	ch := make(chan struct{})
	close(ch)
	ch <- struct{}{}
}
