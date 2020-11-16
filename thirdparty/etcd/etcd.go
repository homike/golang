package etcd

import (
	"bytes"
	"context"
	"fmt"
	"gotest/gologs/logger"
	"io"
	"strconv"
	"strings"
	"time"

	etcd3 "github.com/coreos/etcd/clientv3"
)

func GetServerConfig(key, target string) (io.Reader, error) {
	// get endpoints for register dial address
	client, err := etcd3.New(etcd3.Config{Endpoints: strings.Split(target, ",")})
	if err != nil {
		return nil, fmt.Errorf("grpclb: create etcd3 client failed: %v", err)
	}
	defer client.Close()

	// prefix is the etcd prefix/value to watch
	prefix := fmt.Sprintf("/%s/%s", "etcdname", key)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// query addresses from etcd
	resp, err := client.Get(ctx, prefix, etcd3.WithPrefix())
	if err != nil {
		logger.Fatal("AddWatch Error:%v", key)
	}

	//logger.Error("%v, config %v", prefix, resp)

	return bytes.NewReader(resp.Kvs[0].Value), nil
}

func unmarshalConfig(resp *etcd3.GetResponse) map[int]int {
	serverInfo := map[int]int{}

	if resp == nil || resp.Kvs == nil {
		return serverInfo
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			key := string(resp.Kvs[i].Key)
			sl := strings.Split(key, "/")
			if len(sl) <= 0 {
				logger.Error("Split resp.Kvs[i].Key:%v", key)
				continue
			}
			idStr := sl[len(sl)-1]
			id, err := strconv.Atoi(idStr)
			if nil != err {
				logger.Error("Atoi resp.Kvs[i].Key:%v", key)
				continue
			}
			count, err := strconv.Atoi(string(v))
			if nil != err {
				logger.Error("Atoi resp.Kvs[i].Key:%v", string(v))
				continue
			}

			serverInfo[id] = count
		}
	}

	return serverInfo
}
