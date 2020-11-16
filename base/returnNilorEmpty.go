package base

import (
	"fmt"
	"testing"
)

// nil是首选样式。当我们只需要返回一个空切片nil时，几乎在所有情况下都可以正常工作。
// nil比[]string{}或make([]string, 0)容易键入, 并且通常会突出显示语法，这使它更易于阅读。
func ReturnSlice() []string {
	//return nil
	return []string{}
}

func ReturnMap() map[int]int {
	//return nil
	return make(map[int]int)
}

func Test_Return(t *testing.T) {
	s := []string{}
	ret := ReturnSlice()
	if ret == nil {
		fmt.Println("ret is nil")
	}
	s = append(s, ret...)
	fmt.Printf("slice len: %v, value: %v \n", len(s), s)

}
