package httptest

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func FilteredSQLInject(to_match_str string) bool {
	//fmt.Println("test sql:", to_match_str)
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|((?i)(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b))`
	re, err := regexp.Compile(str)
	if err != nil {
		panic(err.Error())
		return false
	}
	return re.MatchString(to_match_str)
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	openid := r.Form.Get("id")
	if FilteredSQLInject(openid) {
		fmt.Println("error openid")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// for k, v := range r.Form {
	//     fmt.Println("key:", k)
	//     fmt.Println("val:", strings.Join(v, ""))
	// }

	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}
func RunHTTPServer() {
	http.HandleFunc("/", sayhelloName)       //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
