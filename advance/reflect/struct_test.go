package reflecter

import (
	"fmt"
	"reflect"
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

func Test_Members() {
	GetMembers(&sr{})
}
