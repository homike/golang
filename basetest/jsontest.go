package basetest

import (
	"encoding/json"
	"fmt"
	"math"
)

type JsonStruct struct {
	Date int64 `json:"date"`
}

func JsonTest() {
	fmt.Println("9223372036854775807")
	j1 := JsonStruct{
		Date: math.MaxInt64,
	}

	fmt.Println(j1.Date)
	jm, _ := json.Marshal(j1)

	j2 := JsonStruct{}
	json.Unmarshal(jm, &j2)
	fmt.Println(j2.Date)
}
