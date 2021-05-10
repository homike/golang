package etcd_naming

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"gotest/thirdparty/etcd/pb"

	clientv3 "go.etcd.io/etcd/client/v3"
	etcdResolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/reflection"
)

var ETCD_ADDR = "192.168.0.18:2379,192.168.0.18:2381,192.168.0.18:2383"

//var kPrefix = "test/"
var kPrefix = ""

type service struct {
	Port int
}

func (s *service) Ping(ctx context.Context, req *pb.PingReq) (res *pb.PingAck, err error) {
	fmt.Printf("Port: %v in requeset, ping req: [%v] ====> \n", s.Port, req.ReqInfo)

	return &pb.PingAck{
		ReqInfo: req.ReqInfo,
		AckInfo: "ack",
	}, nil
}

func Serve(port int) {
	// 创建gRPC服务器
	grpcSvr := grpc.NewServer()
	// 在gRPC服务端注册服务
	pb.RegisterEtcdServiceServer(grpcSvr, &service{Port: port})
	//在给定的gRPC服务器上注册服务器反射服务
	reflection.Register(grpcSvr)

	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	//serviceName := fmt.Sprintf("Ping")
	//serviceName := fmt.Sprintf("Ping/%v", port)
	serviceName := fmt.Sprintf("%vPing", kPrefix)
	etcdService, err := NewService(NewServiceInfo(addr, serviceName), strings.Split(ETCD_ADDR, ","))
	if err != nil {
		fmt.Printf("failed to NewService: %v", err)
		return
	}
	go etcdService.Start()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	if err := grpcSvr.Serve(lis); err != nil {
		fmt.Println(err)
	}
}

func Client() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(ETCD_ADDR, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("failed to clientv3.New: %v", err)
		return
	}

	r, err := etcdResolver.NewBuilder(cli)
	if err != nil {
		log.Fatalf("etcdResolver.NewBuilder: %v", err)
		return
	}

	target := fmt.Sprintf("%v:///%v%v", r.Scheme(), kPrefix, "Ping")
	//target := fmt.Sprintf("%v:///%v", r.Scheme(), "Ping")
	conn, err := grpc.DialContext(
		context.TODO(),
		target,
		grpc.WithInsecure(),
		grpc.WithResolvers(r),
		// 设置负载均衡策略
		// grpc.RoundRobin()，和grpc.WithBalancer()来设置负载均衡，
		// 这个版本grpc.RoundRobin()已经取消了，grpc.WithBalancer()和grpc. 也WithBalancerName()标记为废弃。现在改为读取外部配置，主要是方便服务启动后动态更新(设计初衷应该是主要用在服务端)
		//grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		return
	}
	fmt.Println("connection target ", conn.Target())

	for i := 0; i < 5; i++ {
		go func(index int) {
			for {
				time.Sleep(500 * time.Millisecond)

				etcdService := pb.NewEtcdServiceClient(conn)
				_, _ = etcdService.Ping(context.TODO(), &pb.PingReq{
					ReqInfo: fmt.Sprintf("req_%v", index),
				})
				//log.Printf("<====Ping Recv Ack: %v, error: %v", resp, err)
			}
		}(i)
	}
}

func TestEtcd(t *testing.T) {
	go Serve(10001)
	go Serve(20001)

	Client()

	for {
	}
}
