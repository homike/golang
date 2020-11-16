package basetest

import (
	"encoding/json"
	"fmt"
	"math"
)

type JsonStruct struct {
	Date int64 `json:"Date"`
}

type JsonStruct2 struct {
	Date int64 `json:"date"`
}

func JsonTest() {
	//fmt.Println("9223372036854775807")
	j1 := JsonStruct{
		Date: math.MaxInt64,
	}

	fmt.Println(j1.Date)
	//jm, _ := json.Marshal(j1)

	strJson := `{"Date": 123}`
	j2 := JsonStruct2{}
	json.Unmarshal([]byte(strJson), &j2)
	fmt.Println(j2.Date)
}
