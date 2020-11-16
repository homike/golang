# gotest
It's userful example for golang api test.

# 目录
    - advance       高级特性
    - algorithm     算法
    - base          基础函数
    - bugs          bug记录
    - design        设计模式 
    - network       网络库
    - thirdparty    第三方库

# cmd

## go test
        -v              显示详细的输出
        -bench regexp   执行相应的 benchmarks, 例如 -bench=.

### test example
    - go test -v helloworld_test.go
    - go test -v -run TestA helloworld_test.go
```
TestHelloWorld(t *testing.T) {
        t.Log("hello world")
}
```

### benchmark example
    - go test -v -bench=. benchmark_test.go
    - 
```
Benchmark_Add(b *testing.B) {
    // 重置计时器
    b.ResetTimer()
    // 停止计时器
    b.StopTimer()
    // 开始计时器
    b.StartTimer()

    var n int
    for i := 0; i < b.N; i++ {
        n++
    }
}
```

## go run 
    - go run -race
    - go build -race
检测竞争
