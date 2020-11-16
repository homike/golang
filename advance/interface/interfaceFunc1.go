package interfacetest

import "fmt"

type HandlerFunc1 func(k, v interface{})

func (f HandlerFunc1) Do(k, v interface{}) {
	f(k, v)
}

type welcome1 string

func (w welcome1) selfInfo(k, v interface{}) {
	fmt.Printf("%s,我叫%s,今年%d岁\n", w, k, v)
}

func RunInterface1() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26

	var w welcome1 = "大家好"

	Each(persons, HandlerFunc1(w.selfInfo))
}
