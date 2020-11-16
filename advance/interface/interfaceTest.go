package interfacetest

import (
	"encoding/json"
	"fmt"
)

/*****************************************************************************************
* 测试1， 直接存储接口, 打印输出为空
*****************************************************************************************/
type Test struct {
	Data map[int]Tester `json:"data"`
}

type Tester interface {
	DoTest()
}

// Father
type Father struct {
	Var1 uint `json:"var1"`
}

func NewFather(v1 uint) *Father {
	return &Father{
		Var1: v1,
	}
}

func (f *Father) DoTest() {
	fmt.Println("value: ", f.Var1)
}

// Child
type Child struct {
	Var1 uint `json:"var1"`
	Var2 uint `json:"var2"`
}

func NewChild(v1, v2 uint) *Child {
	return &Child{
		Var1: v1,
		Var2: v2,
	}
}

func (c *Child) DoTest() {
	fmt.Println("v1: ", c.Var1, " v2:", c.Var2)
}

func Test_OOP() {
	v1 := (interface{})(NewFather(1)).(Tester)
	v2 := (interface{})(NewChild(1, 2)).(Tester)

	m := make(map[int]Tester)
	m[1] = v1
	m[2] = v2
	data1 := Test{
		Data: m,
	}

	jsonData, err := json.Marshal(data1)
	if err != nil {
		return
	}

	fmt.Println("data: ", string(jsonData))
	// data:  {"data":{"1":{"var1":1,"var2":2,"var3":3,"var4":4}, "2":{"var1":1,"var2":2}}}
}

/*****************************************************************************************
* 测试2， 直接存储接口, 打印输出为空
*****************************************************************************************/

type Baseer interface {
	DoTest() uint
	DoTest2() uint
}

type Base2 struct {
	Var1 uint `json:"var1"`
	Var2 uint `json:"var2"`
}

func (b *Base2) DoTest() uint {
	return b.Var1
}

func (b *Base2) DoTest2() uint {
	return b.Var2
}

type Derived2 struct {
	Var3   uint `json:"var3"`
	Var4   uint `json:"var4"`
	*Base2 `json:"father"`
}

func NewDerived2(v1, v2, v3, v4 uint) *Derived2 {
	return &Derived2{
		Base2: &Base2{
			Var1: v1,
			Var2: v2,
		},
		Var3: v3,
		Var4: v4,
	}
}

func (d *Derived2) DoTest3() uint {
	return d.Var3
}

type Test2 struct {
	Data map[int]*Baseer `json:"data"`
}

func Test_OOP2() {
	//Test2
	v1 := new(Baseer)
	*v1 = NewDerived2(1, 2, 3, 4)

	ret1 := (*v1).DoTest()
	ret2 := (*v1).DoTest2()

	c, ok := (interface{})(*v1).(*Child)
	if !ok {

	}
	ret3 := c.DoTest3()
	fmt.Println("resut: ", ret1, ret2, ret3)

	//Test1
	m := make(map[int]*Baseer)
	m[1] = v1

	data1 := Test2{
		Data: m,
	}

	bData, err := json.Marshal(data1)
	if err != nil {
		return
	}

	fmt.Println("data: ", string(bData))
	//{"data":{"1":{}}}
}
