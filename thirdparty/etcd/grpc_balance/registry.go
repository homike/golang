package etcd_naming

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 服务信息
type ServiceInfo struct {
	Addr string
	// 当使用metadata时, 会有报错
	Metadata1 map[string]string
}

const kMetaData_Name = "Name"

func NewServiceInfo(addr, name string) ServiceInfo {
	info := &ServiceInfo{
		Addr:      addr,
		Metadata1: make(map[string]string),
	}
	info.Metadata1[kMetaData_Name] = name

	return *info
}

func (s *ServiceInfo) Name() string {
	return s.Metadata1[kMetaData_Name]
	//return "Ping"
}

func (s *ServiceInfo) IP() string {
	return s.Addr
}

type Service struct {
	ServiceInfo ServiceInfo
	stop        chan error
	leaseId     clientv3.LeaseID
	client      *clientv3.Client
}

// NewService 创建一个注册服务
func NewService(info ServiceInfo, endpoints []string) (service *Service, err error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 5,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	service = &Service{
		ServiceInfo: info,
		client:      client,
	}
	return
}

// Start 注册服务启动
func (service *Service) Start() (err error) {
	ch, err := service.keepAlive()
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		select {
		case err := <-service.stop:
			return err
		case <-service.client.Ctx().Done():
			return errors.New("service closed")
		case resp, ok := <-ch:
			// 监听租约
			if !ok {
				log.Println("keep alive channel closed")
				return service.revoke()
			}
			log.Printf("Recv reply from service: %s, ttl:%d", service.getKey(), resp.TTL)
		}
	}
}

func (service *Service) Stop() {
	service.stop <- nil
}

func (service *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	val, _ := json.Marshal(&service.ServiceInfo)

	// 创建一个租约
	resp, err := service.client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = service.client.Put(context.TODO(), service.getKey(), string(val), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	service.leaseId = resp.ID
	return service.client.KeepAlive(context.TODO(), resp.ID)
}

func (service *Service) revoke() error {
	_, err := service.client.Revoke(context.TODO(), service.leaseId)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", service.getKey())
	return err
}

func (service *Service) getKey() string {
	return service.ServiceInfo.Name() + "/" + service.ServiceInfo.IP()
	//ip := strings.Replace(service.ServiceInfo.IP(), ":", "", -1)
	//return service.ServiceInfo.Name() + "/" + ip
}
