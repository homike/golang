package httptest

import (
	"fmt"
	"reflect"
)

// func GetMembers(i interface{}) {
// 	t := reflect.TypeOf(i)

// 	for {
// 		if t.Kind() == reflect.Struct {
// 			fmt.Printf("%v, %v 个字段\n", t, t.NumField())
// 			for i := 0; i < t.NumField(); i++ {
// 				fmt.Println(t.Field(i).Name)
// 			}
// 		}
// 		fmt.Printf("%v, %v 个字段\n", t, t.NumMethod())

// 		for i := 0; i < t.NumMethod(); i++ {
// 			fmt.Println(t.Method(i).Name)
// 		}

// 		if t.Kind() == reflect.Ptr {
// 			t = t.Elem()
// 		} else {
// 			break
// 		}
// 	}
// }

// type sr struct {
// 	string
// }

// func (s sr) Read() {
// }

// func (s *sr) Write() {
// }

// func main() {
// 	GetMembers(&sr{})
// }

///////////////////////////////
// Get Struct Tag
//////////////////////////////
// func main() {
// 	type S struct {
// 		F string `name:"test" color:"blue"`
// 	}

// 	s := S{}
// 	st := reflect.TypeOf(s)
// 	for i := 0; i < st.NumField(); i++ {
// 		fmt.Println(st.Field(i).Tag.Get("name"), st.Field(i).Tag.Get("color"))
// 	}
// }

///////////////////////////////
// Get Struct Tag
//////////////////////////////
func RunRelect() {
	type S struct {
		S1 struct {
			S1F string `name:"test" color:"blue"`
		} `id:"1"`

		S2 struct {
			S2F string `name:"test" color:"green"`
		} `id:"2"`
	}

	s := S{}
	sVal := reflect.ValueOf(&s).Elem()
	sType := sVal.Type()

	for i := 0; i < sVal.NumField(); i++ {
		vf := sVal.Field(i)
		vfT := vf.Type()

		tf := sType.Field(i)
		fmt.Println("id", tf.Tag.Get("id"))

		for j := 0; j < vf.NumField(); j++ {
			fmt.Println("--color", vfT.Field(j).Tag.Get("color"))
		}
	}
}

// type T struct {
// 	A int    `name:"czx" color:"A"`
// 	B string `name:"czx" color:"B"`
// }

// func main() {
// 	t := T{12, "skidoo"}

// 	sVal := reflect.ValueOf(&t).Elem()
// 	typeOfT := sVal.Type()

// 	for i := 0; i < sVal.NumField(); i++ {
// 		vf := sVal.Field(i)
// 		tf := typeOfT.Field(i)
// 		fmt.Println(tf.Name, vf.Type(), vf.Interface(), tf.Tag.Get("name"))
// 	}
// }
