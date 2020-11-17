// nil_test.go 是来显示golang interface nil的坑的,
// golang interface内部实现如下, 这就导致将要给空结构体指针，赋值给interface时, 虽然值内容是nil, 但是值类型不是nil, 因此导致无法判断==nil
//type Interface struct {
//	pt uintptr // 到值类型的指针
//	pv uintptr // 到值内容的指针
//}

package interfacetest

import (
	"fmt"
	"reflect"
	"testing"
)

type NilTest struct {
}

var niltest *NilTest

// 返回特性指针类型(output1)是不会有问题的, 但是返回interface{} (output2)时, 判断nil时会出现问题
// `避免将一个有可能为 nil 的具体类型的值赋值给 interface 变量`
func output1() interface{} {
	return niltest
}

// 解决方案
// 1. 返回特定类型的指针, 不要返回interface{}
func output2() *NilTest {
	return niltest
}

// 2. 返回前判断指针是否为nil
func output3() interface{} {
	if niltest == nil {
		return nil
	}
	return niltest
}

// 3. 通过IsNil() 来判断
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func TestNil(t *testing.T) {
	o1, o2, o3 := output1(), output2(), output3()

	fmt.Println(o1 == nil, o2 == nil, o3 == nil)

	fmt.Println(IsNil(o1), IsNil(o2), IsNil(o3))
}
