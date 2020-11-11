package handler

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

var (
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfBytes   = reflect.TypeOf(([]byte)(nil))
	typeOfSession = reflect.TypeOf(&Session{})
)

func isExported(name string) bool {
	w, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(w)
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		fmt.Println("000")
		return false
	}

	// Method needs three ins: receiver, *Session, []byte or pointer.
	if mt.NumIn() != 3 {
		fmt.Println("111")
		return false
	}

	// Method needs one outs: error
	if mt.NumOut() != 1 {
		fmt.Println("222")
		return false
	}

	if t1 := mt.In(1); t1.Kind() != reflect.Ptr || t1 != typeOfSession {
		fmt.Println("333")
		return false
	}

	if (mt.In(2).Kind() != reflect.Ptr && mt.In(2) != typeOfBytes) || mt.Out(0) != typeOfError {
		fmt.Println("444")
		return false
	}
	return true
}
