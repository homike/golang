package interfacer

import "fmt"

type peopleer interface {
	eat()
	say()
}
type father struct {
	cname string
}

func (f *father) eat() {
	fmt.Println(f.cname, "eat")
}
func (f *father) say() {
	fmt.Println(f.cname, "say")
}

type child struct {
	father
	cname string
}

func (c *child) eat() {
	//c.father.eat()
	fmt.Println(c.cname, "eat1")
}

// func (c *child) say() {
// 	//c.father.say()
// 	fmt.Println(c.cname, "say1")
// }

func factory() interface{} {
	c := &child{cname: "child"}
	return c
}

func Run1() {
	ret := factory()
	f, _ := ret.(peopleer)
	//fmt.Println(ok)
	f.eat()

	fmt.Println("-----------------------")

	obj := child{
		cname: "child",
		father: father{
			cname: "father",
		},
	}
	var people peopleer = &obj
	people.eat()
	people.say()

	var fpeople peopleer = &obj.father
	fpeople.eat()
	fpeople.say()
}
