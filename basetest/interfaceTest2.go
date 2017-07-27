package basetest

import (
	"fmt"
)

type commonIF interface {
	GetVal1() int
	//GetVal2() int
}

type Base struct {
	Val1 int
	Val2 int
}

func (b *Base) GetVal1() int {
	return b.Val1
}

func (b *Base) GetVal2() int {
	return b.Val2
}

//--------------------------------------------------------------
type Foo struct {
	Val3 int
	*Base
}

func (f *Foo) GetVal3() int {
	return f.Val3
}

func TestIF2() {
	comIF := new(commonIF)
	(*comIF) = &Foo{
		Base: &Base{
			Val1: 1,
			Val2: 2,
		},
		Val3: 3,
	}
	fmt.Println((*comIF).GetVal1())

	foo, ok := (*comIF).(*Foo)
	if !ok {
		fmt.Println("Error Convert")
		return
	}
	fmt.Println(foo.GetVal3())

	// base, ok := (*comIF).(*Base)
	// if !ok {
	// 	fmt.Println("Error Convert")
	// 	return
	// }
	// fmt.Println(base.GetVal2())

	//FrameWork(foo.Base)
	//fmt.Println((*comIF).GetVal2())
}

// //--------------------------------------------------------------
// type commonIF interface {
// 	GetVal1() int
// 	//GetVal2() int
// }

// type Base struct {
// 	Val1 int
// 	Val2 int
// }

// func (b *Base) GetVal1() int {
// 	return b.Val1
// }

// func (b *Base) GetVal2() int {
// 	return b.Val2
// }

// func FrameWork(cf commonIF) {
// 	base, ok := (cf).(*Base)
// 	if !ok {
// 		fmt.Println("Error Convert")
// 		return
// 	}
// 	fmt.Println(base.GetVal2())
// }

// //--------------------------------------------------------------
// type Foo struct {
// 	Val3 int
// 	commonIF
// }

// func (f *Foo) GetVal3() int {
// 	return f.Val3
// }

// func TestIF2() {
// 	comIF := new(commonIF)
// 	(*comIF) = &Foo{
// 		commonIF: &Base{
// 			Val1: 1,
// 			Val2: 2,
// 		},
// 		Val3: 3,
// 	}
// 	fmt.Println((*comIF).GetVal1())

// 	foo, ok := (*comIF).(*Foo)
// 	if !ok {
// 		fmt.Println("Error Convert")
// 		return
// 	}
// 	fmt.Println(foo.GetVal3())

// 	FrameWork(foo.commonIF)
// }
