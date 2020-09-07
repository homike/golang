package inject

import (
	"fmt"

	wire "github.com/google/wire"
)

//  google 的方式是通过生成代码的方式，来解决依赖，这样就不用等待运行时才出错
type Foo struct {
}

func NewFoo() *Foo {
	return &Foo{}
}

type Bar struct {
	foo *Foo
}

func NewBar(foo *Foo) *Bar {
	return &Bar{
		foo: foo,
	}
}

func (p *Bar) Test() {
	fmt.Println("hello")
}

//---------------------------------------------------
type Instance struct {
	Foo *Foo
	Bar *Bar
}

var SuperSet = wire.NewSet(NewFoo, NewBar)

func InitializeAllInstance() *Instance {
	wire.Build(SuperSet, Instance{})
	return &Instance{}
}
