# ETCD + gRPC 服务发现

### 目录
- custom_naming
    自定义resolver.Builder

- etcd_naming
    使用etcd提供的go.etcd.io/etcd/client/v3/naming/resolver，来完成resolver.Builder接口
    但是etcd中的数据需要遵循如下数据格式
    ```
    EtcdDat struct {
        Addr string
        Metadata interface{}
    }
    ```

### gRPC Name Resolution
它定义了 gRPC 的 URI 格式，举个例子 「dns://1.1.1.1/huoding.com」，其中： 
- Scheme：dns
- Authority：1.1.1.1
- Endpoint：huoding.com
表示通过 dns 服务器 1.1.1.1 查询 huoding.com 有哪些节点

既然我们要支持 etcd，那么我们首先要想好 etcd 对应的 URI 应该是什么样的，Authority 填什么好呢？
按 dns 的例子的意思，填 etcd 服务器的地址似乎就可以，不过实际情况中，一般会有多台 etcd 服务器，还牵扯到用户名密码，与其全写到 Authority 里，
还不如直接从配置文件里获取来的实在，所以建议 Authority 留空。假设我们要通过 etcd 查询一个名为 foo 的服务对应的节点，那么 URI 可以定义为：
**「etcd:///foo」**

### gRPC + ETCD 服务发现流程
- 服务端启动，在 etcd 里通过租约注册键为「/foo/<ip>:<port>」并且值为「<ip>:<port>」的数据，同时定期发送心跳包，一旦节点退出会注销相关数据

- 客户端启动，gRPC 把 etcd:///foo 解析出 Scheme、Authority、Endpoint，并根据 Scheme 找到对应 Builder，调用其 Build 方法，返回对应的 Resolver，
在 etcd 中查询前缀是「/foo/」的数据，就是目前可用的节点。

- 最后，负载均衡会挑选出一个节点来提供服务。


