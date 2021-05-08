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
	"google.golang.org/grpc/reflection"
)

var ETCD_ADDR = "192.168.0.18:2379,192.168.0.18:2381,192.168.0.18:2383"

type service struct {
}

func (s *service) Ping(ctx context.Context, req *pb.PingReq) (res *pb.PingAck, err error) {
	fmt.Printf("ping ack: [%v] ====> \n", req.ReqInfo)
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

	etcdService, err := NewService(NewServiceInfo("127.0.0.1:8972", "Ping"), strings.Split(ETCD_ADDR, ","))
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

	watchTarget := fmt.Sprintf("%v:///%v", r.Scheme(), "Ping")
	conn, err := grpc.DialContext(context.TODO(), watchTarget, grpc.WithInsecure(), grpc.WithResolvers(r))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		return
	}
	fmt.Println("target ", conn.Target())

	etcdService := pb.NewEtcdServiceClient(conn)
	resp, err := etcdService.Ping(context.TODO(), &pb.PingReq{
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
