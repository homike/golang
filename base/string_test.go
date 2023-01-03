package base

import (
	"fmt"
	"testing"
)

type Id struct {
	I int64
}

func (i *Id) String() string {
	return fmt.Sprintf("uid:%d", i.I)
}

var i = &Id{
	I: 64,
}

func Benchmark_String(b *testing.B) {
	fmt.Printf("%s", i)
}

func Benchmark_NoString(b *testing.B) {
	fmt.Printf("uid:%d", i.I)
}
