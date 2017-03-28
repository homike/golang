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
	req, _ := http.NewRequest("POST", "http://http://127.0.0.1:8680/x/3/idip", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req) //发送
	_ = err
	return resp
}

func RunHTTPClient() {
	forbinUser := `data_packet={"head":{"PacketLen":12,"Cmdid":4157,"Seqid":123,"ServiceName":"idip", "SendTime": "19990202", "Version":123, "Authenticate": "ffffefe","Result":-1,"RetErrMsg":"fff"},"body":{"AreaId":999,"PlatId":1,"OpenId":"B8AD3948B96CCEF3547E51D728C37B09","Type":2,"msg":"别玩啦，你妈喊你回家吃饭","Sender":"2100","Source":11,"Serial":"Serial"}}`
	r := HttpPost(forbinUser)
	var x map[string]interface{}

	_ = r
	_ = x
	//err := json.Unmarshal(data, v)
	//fmt.Println(err)
	//fmt.Printf("%+v", x)
}
