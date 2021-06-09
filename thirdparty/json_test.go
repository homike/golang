package thirdparty

import (
	"encoding/json"
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestJson(t *testing.T) {
	type Message struct {
		Name string `json:"msg_name"`       // 对应JSON的msg_name
		Body string `json:"body,omitempty"` // 如果为空置则忽略字段
		Time int64  `json:"-"`              // 直接忽略字段
	}
	var m = Message{
		Name: "Alice",
		Body: "",
		Time: 1294706395881547000,
	}
	data, err := json.Marshal(m)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	//fmt.Println(string(data))

	//Output:
	//{"msg_name":"Alice"}
	type MessageLow struct {
		name string `json:"msg_name"`       // 对应JSON的msg_name
		body string `json:"body,omitempty"` // 如果为空置则忽略字段
		time int64  `json:"-"`              // 直接忽略字段
	}
	decodeJson := MessageLow{}
	_ = json.Unmarshal(data, &decodeJson)
	fmt.Println("decode", decodeJson)
	//Output: decode {Alice  0}
}

func TestJson2(t *testing.T) {
	var JsonIter = jsoniter.ConfigCompatibleWithStandardLibrary

	type Message struct {
		Name string `json:"msg_name"`       // 对应JSON的msg_name
		Body string `json:"body,omitempty"` // 如果为空置则忽略字段
		Time int64  `json:"-"`              // 直接忽略字段
	}
	var m = Message{
		Name: "Alice",
		Body: "",
		Time: 1294706395881547000,
	}
	data, err := JsonIter.Marshal(m)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	//fmt.Println(string(data))

	//Output:
	//{"msg_name":"Alice"}
	type MessageLow struct {
		name string `json:"msg_name"`       // 对应JSON的msg_name
		body string `json:"body,omitempty"` // 如果为空置则忽略字段
		time int64  `json:"-"`              // 直接忽略字段
	}
	decodeJson := MessageLow{}
	_ = JsonIter.Unmarshal(data, &decodeJson)
	fmt.Println("decode", decodeJson)
	//Output: decode {Alice  0}
}
