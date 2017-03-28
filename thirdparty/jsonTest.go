package thirdparty

import (
	"encoding/json"
	"fmt"
)

type ColorGroup struct {
	ID     int
	Name   string
	Colors []string
}

func main() {

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
	fmt.Println(string(data))

	//Output:
	//{"msg_name":"Alice"}
	decodeJson := Message{}
	_ = json.Unmarshal(data, &decodeJson)
	fmt.Println("decode", decodeJson)
	//Output: decode {Alice  0}
}
