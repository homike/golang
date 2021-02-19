package bugs

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

type MyContext struct {
	context.Context
	v int
}

func blockByContext(ctx *MyContext, timeout time.Duration) {
	hctx, hcancel := context.WithCancel(ctx)
	if timeout > 0 {
		hctx, hcancel = context.WithTimeout(ctx, timeout)
	}
	_, _ = hctx, hcancel
}

// 02057906f7272a4787b8a0b5b7cafff8ad3024f0
// 这一次提交之前, 如果将context封装到自定义结构体中
// 在调用withCancel, withTimeout时, 将会启动意外的gorutine
// 目前的版本已经修复, 所有该bug也就不存在了
func TestBlockByContext(t *testing.T) {
	fmt.Println("Start Gorutines: ", runtime.NumGoroutine())

	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	myCtx := &MyContext{
		Context: ctx,
		v:       1,
	}

	go blockByContext(myCtx, time.Second)

	fmt.Println("End Gorutines: ", runtime.NumGoroutine())
	for {
	}
}
