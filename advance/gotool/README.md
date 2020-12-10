# gotool
    gotool 是golang的性能诊断工具用例, 包含pprof, trace, 火焰图等信息

# 目录
    - trace    xxxx 

# 安装
    - 支持svg矢量图 
    安装[Graphviz](https://graphviz.org/download/)

    - 支持火焰图
    安装FlameGraph
    ```
        git clone https://github.com/brendangregg/FlameGraph.git
        cp flamegraph.pl /usr/local/bin
    ```
    安装go-torch
    ```
        go get -v github.com/uber/go-torch
    ```

    - 支持trace view
    如果trace view 是空白, 是golang版本的问题, 可以安装[gotip](https://godoc.org/golang.org/dl/gotip)
        然后用gotip 代替 go 来运行, 如:  gotip tool xxxx

# cmd

## trace
    
### go tool
        -http           指定http展示addr 

### example
    - 通过程序中已经开启的pprof接口, 来获取trace信息
    ```
    curl -XGET "http://127.0.0.1:6060/debug/pprof/trace?seconds=30" -o trace.out
    ```
    - 打开web页面, 进行数据分析
    go/gotip tool trace -http=0.0.0.0:12000 trace.out 

```
package main
import (	
    "os"
    "runtime/trace"
)

func main() {
    f, err := os.Create("trace.out")	
    if err != nil {		
       panic(err)
    }	
    defer f.Close()

    err = trace.Start(f)
     if err != nil {
 	panic(err)
    }	
    defer trace.Stop()  
    // Your program here
}
```

## 火焰图

    go-torch -u http://localhost:9090 -t 30
