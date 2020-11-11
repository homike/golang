package methodhandler

import (
	"fmt"
	"gotest/gologs/logger"
	"reflect"
	"strings"
)

const (
	METHOD_TEST = "test"
)

var (
	__callMap = map[string]*CallFunc{}
)

type CallFunc struct {
	Func       interface{}
	FuncType   reflect.Type
	FuncVal    reflect.Value
	FuncParams string
}

func findCall(funcName string) *CallFunc {
	funcName = strings.ToLower(funcName)

	fun, exist := __callMap[funcName]
	if exist == true {
		return fun
	}
	return nil
}

func RegisterCall(funcName string, call interface{}) {
	funcName = strings.ToLower(funcName)
	if findCall(funcName) != nil {
		logger.Error("Call already register")
		return
	}

	__callMap[funcName] = &CallFunc{Func: call, FuncVal: reflect.ValueOf(call), FuncType: reflect.TypeOf(call), FuncParams: reflect.TypeOf(call).String()}
}

// params 中不能有nil, 否则无法执行
func ExecFunc(funcName string, params []interface{}) (interface{}, interface{}) {
	pFunc := findCall(funcName)
	if pFunc != nil {
		f := pFunc.FuncVal
		k := pFunc.FuncType
		strParams := pFunc.FuncParams

		if k.NumIn() != len(params) {
			logger.Error("func [%s] can not call, func params [%s], params [%v]", funcName, strParams, params)
			return nil, fmt.Errorf("params not enough")
		}

		if len(params) >= 1 {
			in := make([]reflect.Value, len(params))
			for i, param := range params {
				v := reflect.ValueOf(param)
				if v.IsNil() {
					logger.Error("func [%s] params has nil", funcName)
					return nil, fmt.Errorf("params has nil")
				}
				in[i] = v
			}

			//in 参数中如果有nil, 会panic
			ret := f.Call(in)
			if len(ret) >= 2 && ret[0].CanInterface() && ret[1].CanInterface() {
				return ret[0].Interface(), ret[1].Interface()
			}
		} else {
			logger.Error("func [%s] params at least one", funcName)
			return nil, fmt.Errorf("params not enough")
		}
	}

	return nil, fmt.Errorf("execfunc failed")
}
