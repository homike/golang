package base

import (
	"fmt"
	"testing"
	"unicode"
)

func Test_Han(t *testing.T) {
	name := "我的测试字符串12,@"
	right := true
	for _, v := range name {
		if unicode.Is(unicode.Scripts["Han"], v) || unicode.IsLetter(v) || unicode.IsNumber(v) {
			continue
		}
		right = false
		break
	}
	fmt.Println("is right name: ", right)
}

func Benchmark_Han(b *testing.B) {
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	b.StopTimer()

	name := "我的测试字符串12,@"

	// 开始计时器
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range name {
			if unicode.Is(unicode.Scripts["Han"], v) || unicode.IsLetter(v) || unicode.IsNumber(v) {
				continue
			}
			break
		}
	}
}
