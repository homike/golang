package httptest

import (
	"fmt"
	"log"
	"net/http"
)

////////////////////////////////////
// golang 使用web的调用方式
///////////////////////////////////

// 方法一: 使用默认mux
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)

	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func RunHTTPServer1() {
	http.HandleFunc("/", sayhelloName)       //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 方法二: 自定义mux
type a struct {
}

func (*a) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()
	switch path {
	case "/":
		fmt.Println("a.ServeHTTP /")
	case "/test":
		fmt.Println("a.ServeHTTP /test")
	}
}

func RunHTTPServer2() {
	err := http.ListenAndServe(":9090", &a{}) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 方法三: 使用系统提供的 mux
func RunHTTPServer3() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", sayhelloName)
	mux.HandleFunc("/test", sayhelloName)

	err := http.ListenAndServe(":9090", mux) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
