package interfacetest

import (
	"fmt"
	"reflect"
)

func TestNil() {
	var a interface{} = nil // tab = nil, data = nil
	var b interface{} =

	//(*int)(nil) // tab 包含 *int 类型信息, data = nil

	fmt.Println(IsNil(a))

	fmt.Println(IsNil(b))
}

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
