package base

import (
	"fmt"
	"testing"
)

func Test_Bit(t *testing.T) {
	//var tnum uint32 = 7
	var tdata uint32 = 0

	// 记录某一位是否被操作
	tdata |= 1 << 7
	//tdata |= 1 << 8

	// 判断某一个是否被操作
	ret := (tdata >> 7) & 1 // 1/0

	fmt.Println("tdata:", tdata)
	fmt.Println("result: ", ret)
	fmt.Println("result: ", (tdata>>8)&1)
	fmt.Println("result: ", (tdata>>9)&1)

}
