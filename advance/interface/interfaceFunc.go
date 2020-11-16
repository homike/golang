package interfacetest

import "fmt"

type HandlerInterface interface {
	Do(k, v interface{})
}

func Each(m map[interface{}]interface{}, h HandlerInterface) {
	if m != nil && len(m) > 0 {
		for k, v := range m {
			h.Do(k, v)
		}
	}
}

type welcome struct {
	name string
}

func (w welcome) Do(k, v interface{}) {
	fmt.Printf("%s,我叫%s,今年%d岁\n", w.name, k, v)
}

func Run() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26

	var w welcome = welcome{
		name: "大家好",
	}

	Each(persons, w)
}
