package etcd

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"gotest/thirdparty/etcd/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/resolver"
)

var ETCD_ADDR = "192.168.0.18:2379,192.168.0.18:2381,192.168.0.18:2383"

type service struct {
}

func (s *service) Ping(ctx context.Context, req *pb.PingReq) (res *pb.PingAck, err error) {
	fmt.Printf("ping ack: %v ====> \n", req.ReqInfo)
	return &pb.PingAck{
		ReqInfo: req.ReqInfo,
		AckInfo: "ack",
	}, nil
}

func Serve() {
	// 创建gRPC服务器
	grpcSvr := grpc.NewServer()
	// 在gRPC服务端注册服务
	pb.RegisterEtcdServiceServer(grpcSvr, &service{})
	//在给定的gRPC服务器上注册服务器反射服务
	reflection.Register(grpcSvr)

	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。

	etcdService, err := NewService(ServiceInfo{Name: "Ping", IP: "127.0.0.1:8972"}, strings.Split(ETCD_ADDR, ","))
	if err != nil {
		fmt.Printf("failed to NewService: %v", err)
		return
	}
	go etcdService.Start()

	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	if err := grpcSvr.Serve(lis); err != nil {
		fmt.Println(err)
	}
}

func Client() {
	r := NewResolver(strings.Split(ETCD_ADDR, ","), "Ping")
	resolver.Register(r)

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	// The gRPC client library will use the specified scheme to pick the right resolver plugin and pass it the fully qualified name string.

	addr := fmt.Sprintf("%s:///%s", r.Scheme(), "Ping" /*测试，这个可以随便写，底层只是取scheme对应的Build对象*/)

	log.Print("addr: ", addr)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(),
		// grpc.WithBalancerName(roundrobin.Name),
		//指定初始化round_robin => balancer (后续可以自行定制balancer和 register、resolver 同样的方式)
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithBlock())

	// 这种方式也行
	//conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	//conn, err := grpc.Dial(":8972", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	/*conn, err := grpc.Dial(
	      fmt.Sprintf("%s://%s/%s", "consul", GetConsulHost(), s.Name),
	      //不能block => blockkingPicker打开，在调用轮询时picker_wrapper => picker时若block则不进行robin操作直接返回失败
	      //grpc.WithBlock(),
	      grpc.WithInsecure(),
	      //指定初始化round_robin => balancer (后续可以自行定制balancer和 register、resolver 同样的方式)
	      grpc.WithBalancerName(roundrobin.Name),
	      //grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	  )
	  //原文链接：https://blog.csdn.net/qq_35916684/article/details/104055246*/

	if err != nil {
		panic(err)
	}

	client := pb.NewEtcdServiceClient(conn)
	resp, err := client.Ping(context.TODO(), &pb.PingReq{
		ReqInfo: "req",
	})
	log.Print("<====Ping Recv Ack", resp)
}

func TestEtcd(t *testing.T) {
	go Serve()

	Client()

	for {
	}
}
