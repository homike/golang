module gotest

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/facebookarchive/inject v0.0.0-20180706035515-f23751cae28b
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/inject v0.0.0-20180706035515-f23751cae28b // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/structtag v0.0.0-20150214074306-217e25fb9691 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/garyburd/redigo v1.6.2
	github.com/go-sql-driver/mysql v1.5.0
	//github.com/golang/protobuf v1.5.2
	github.com/golang/protobuf v1.4.2
	github.com/google/wire v0.4.0
	github.com/gotomicro/ego v0.5.7
	github.com/gotomicro/ego-component/egorm v0.2.1
	github.com/json-iterator/go v1.1.10
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/pkg/errors v0.9.1
	github.com/powerman/rpc-codec v1.2.2
	github.com/stretchr/testify v1.7.0
	gitlab.papegames.com/daozhang/chanrpc v0.1.4
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0
	go.uber.org/atomic v1.9.0
	google.golang.org/grpc v1.32.0
	gopkg.in/square/go-jose.v2 v2.5.1
)

replace github.com/golang/protobuf v1.4.2 => github.com/golang/protobuf v1.3.5

//replace go.etcd.io/bbolt => github.com/coreos/bbolt v1.3.5
//replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
//replace google.golang.org/grpc v1.32.0 => google.golang.org/grpc v1.26.0
