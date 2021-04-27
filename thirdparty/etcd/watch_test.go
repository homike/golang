package etcd

//func addWatch(name, target string) (*grpc.ClientConn, error) {
//	cli, err := clientv3.New(clientv3.Config{
//		Endpoints:   strings.Split(target, ","),
//		DialTimeout: 5 * time.Second,
//	})
//
//	etcdResolver, err := resolver.NewBuilder(cli)
//	if err != nil {
//		return err
//	}
//
//	//grpc target format (etcd://authority/prefix)
//	watchTarget := fmt.Sprintf("etcd:///%v%v", "x5test", name)
//
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	conn, err := grpc.DialContext(ctx, watchTarget, grpc.WithInsecure(), grpc.WithResolvers(etcdResolver))
//	if err != nil {
//		fmt.Printf("error:%v", err)
//		return nil, err
//	}
//
//	return conn, nil
//}
//
//func addWatch(name, target string) (*grpc.ClientConn, error) {
//	cli, err := clientv3.New(clientv3.Config{
//		Endpoints:   strings.Split(target, ","),
//		DialTimeout: 5 * time.Second,
//	})
//
//	etcdResolver, err := resolver.NewBuilder(cli)
//	if err != nil {
//		return err
//	}
//
//	//grpc target format (etcd://authority/prefix)
//	watchTarget := fmt.Sprintf("etcd:///%v%v", "x5test", name)
//
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	conn, err := grpc.DialContext(ctx, watchTarget, grpc.WithInsecure(), grpc.WithResolvers(etcdResolver))
//	if err != nil {
//		fmt.Printf("error:%v", err)
//		return nil, err
//	}
//
//	return conn, nil
//}

/*
func TestWatch(t *testing.T) {
	conn, err := addWatch("watch", "http://192.168.0.18:2379,http://192.168.0.18:2381,http://192.168.0.18:2383")
	if err != nil {
		fmt.Println("addwatch error")
	}
}
*/
