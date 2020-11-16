//  http handler func 的优化过程
package interfacetest

import "fmt"

// 基本的接口的用法, 实现Do接口
type handler interface {
	Do(k, v interface{})
}

func Each(m map[interface{}]interface{}, h handler) {
	if m != nil && len(m) > 0 {
		for k, v := range m {
			h.Do(k, v)
		}
	}
}

type welcome struct {
}

func (w welcome) Do(k, v interface{}) {
	fmt.Printf("我叫%s,今年%d岁\n", k, v)
}

func Test_FuncHandler1() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26

	var w welcome = welcome{}

	Each(persons, w)
}

//-----------------------------------------------------------------------------------------------------------------
// 优化1
// 标准库中HTTP Handler 的用法, 将func作为一个type
// 这样各种HTTP handler可以定义自己的名字, 只要有相同的参数 和 返回值

type HandlerFunc func(k, v interface{})

func (f HandlerFunc) Do(k, v interface{}) {
	f(k, v)
}

func SelfInfo(k, v interface{}) {
	fmt.Printf("大家好, 我叫%s,今年%d岁\n", k, v)
}

func Test_FuncHandler2() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26

	Each(persons, HandlerFunc(SelfInfo))
}

//-----------------------------------------------------------------------------------------------------------------
// 优化2
// 添加EachFunc方法，这样就不需要在每次调用的时候都强制转化HandlerFunc
func EachFunc(m map[interface{}]interface{}, f func(k, v interface{})) {
	Each(m, HandlerFunc(f))
}

func Test_FuncHandler3() {
	persons := make(map[interface{}]interface{})
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26

	EachFunc(persons, SelfInfo)
}
