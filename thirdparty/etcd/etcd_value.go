package etcd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetServerConfig(key, target string) (io.Reader, error) {
	// get endpoints for register dial address
	client, err := clientv3.New(clientv3.Config{Endpoints: strings.Split(target, ",")})
	if err != nil {
		return nil, fmt.Errorf("grpclb: create clientv3 client failed: %v", err)
	}
	defer client.Close()

	// prefix is the etcd prefix/value to watch
	prefix := fmt.Sprintf("/%s/%s", "etcdname", key)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// query addresses from etcd
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		//logger.Fatal("AddWatch Error:%v", key)
	}

	//logger.Error("%v, config %v", prefix, resp)

	return bytes.NewReader(resp.Kvs[0].Value), nil
}

func unmarshalConfig(resp *clientv3.GetResponse) map[int]int {
	serverInfo := map[int]int{}

	if resp == nil || resp.Kvs == nil {
		return serverInfo
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			key := string(resp.Kvs[i].Key)
			sl := strings.Split(key, "/")
			if len(sl) <= 0 {
				//logger.Error("Split resp.Kvs[i].Key:%v", key)
				continue
			}
			idStr := sl[len(sl)-1]
			id, err := strconv.Atoi(idStr)
			if nil != err {
				//logger.Error("Atoi resp.Kvs[i].Key:%v", key)
				continue
			}
			count, err := strconv.Atoi(string(v))
			if nil != err {
				//logger.Error("Atoi resp.Kvs[i].Key:%v", string(v))
				continue
			}

			serverInfo[id] = count
		}
	}

	return serverInfo
}
