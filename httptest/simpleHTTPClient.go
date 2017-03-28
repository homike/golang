package httptest

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type U struct {
	Name string
	Age  int `json:"appid"`
	Sex  string
}

//Post
func HttpPost(json string) *http.Response {

	body := ioutil.NopCloser(strings.NewReader(json)) //把form数据编下码
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://http://127.0.0.1:8680", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req) //发送
	_ = err
	return resp
}

func RunHTTPClient() {
	forbinUser := ``
	r := HttpPost(forbinUser)
	var x map[string]interface{}

	_ = r
	_ = x
	//err := json.Unmarshal(data, v)
	//fmt.Println(err)
	//fmt.Printf("%+v", x)
}
