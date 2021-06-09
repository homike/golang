package reflecter

import (
	"fmt"
	"reflect"
	"testing"
)

func GetMembers(i interface{}) {
	t := reflect.TypeOf(i)

	for {
		if t.Kind() == reflect.Struct {
			fmt.Printf("%v, %v 个字段\n", t, t.NumField())
			for i := 0; i < t.NumField(); i++ {
				fmt.Println(t.Field(i).Name)
			}
		}
		fmt.Printf("%v, %v 个字段\n", t, t.NumMethod())

		for i := 0; i < t.NumMethod(); i++ {
			fmt.Println(t.Method(i).Name)
		}

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		} else {
			break
		}
	}
}

type sr struct {
	string
}

func (s sr) Read() {
}

func (s *sr) Write() {
}

func _Test_Members(t *testing.T) {
	GetMembers(&sr{})
}

type Tester interface {
	Test()
}

type A struct {
}

func (a *A) Test() {
	fmt.Println("aaa")
}

type B struct {
}

func (a *B) Test() {
	fmt.Println("bbb")
}

type C struct {
}

func (a *C) Test() {
	fmt.Println("ccc")
}

type MemberInterface struct {
	Num   int
	Str   string
	AData *A
	BData *B
	CData *C
}

// 测试所有成员实现了特定接口
func TestStructMembersInterface(t *testing.T) {
	m := &MemberInterface{
		AData: &A{},
		BData: &B{},
		CData: &C{},
	}

	ty := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	for k := 0; k < ty.NumField(); k++ {
		vf := v.Field(k)
		if vf.Kind() == reflect.Ptr {
			if app, ok := vf.Interface().(Tester); ok {
				app.Test()
			}
		} else {
			continue
		}
	}
}
