package thirdparty

import (
	"fmt"
	"strings"
	"testing"

	"github.com/garyburd/redigo/redis"
)

var (
	RedisPool *redis.Pool
)

func init() {
	RedisPool = &redis.Pool{
		MaxIdle:     1,
		IdleTimeout: 1,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			//if _, err := c.Do("AUTH", "123456"); err != nil {
			//c.Close()
			//return nil, err
			//}
			return c, nil
		},
	}
	fmt.Println("Redis Init Success")
}

func GetPool() *redis.Pool {
	return RedisPool
}

/*
type AccountBasicInfo struct {
	AccountID  int32   `redis:"accountid"`
	Permission int32   `redis:"permission" default:"1"`
	Token      string  `redis:"token"`
	OpenID     string  `redis:"open_id"`
	Vps        []int64 `redis:"vps"`
}
*/

// 测试redigo/redis 序列化数据
func _TestRedisHgetAll(t *testing.T) {
	con := GetPool().Get()

	data := &AccountBasicInfo{
		OpenID: "123",
	}
	args := redis.Args{}.Add("Test").AddFlat(data)
	_, err := con.Do("HMSET", args...)
	if err != nil {
	}

	value, err := redis.Values(con.Do("HGETALL", "Test"))
	//value, err := redis.Values(con.Do("get", key))
	if err != redis.ErrNil {
	} else {
	}
	if err != nil {
		fmt.Println("hgetall error :", err)
		return
	}
	fmt.Println("value ", value)

	account := &AccountBasicInfo{}
	if err := redis.ScanStruct(value, account); err != nil {
		fmt.Println(err)
	}

	fmt.Println(account)
}

// 测试redis中的数据为空时，是否会返回error
func _TestRedisEmpty(t *testing.T) {
	con := GetPool().Get()

	_, err := redis.Values(con.Do("HGETALL", "xxx"))
	//value, err := redis.Values(con.Do("get", key))
	if err != redis.ErrNil {
		fmt.Println("redis get error: ErrNil")
	} else {
		fmt.Println("redis get error: ", err)
	}
}

type AccountBindInfo struct {
	OpenID    string `xorm:"openid unique pk" redis:"openid,omitempty"`
	AccountID int32  `xorm:"accountid autoincr" redis:"accountid,omitempty"`
	OpenKey   string `xorm:"-" redis:"openkey,omitempty"` // SDK的openKey，在测试创角登录环境下为账号的密码， 不存于mysql
}

func TestCopyAccountBindInfo(t *testing.T) {
	con := GetPool().Get()
	defer con.Close()

	//keys, err := redis.Strings(con.Do("keys", "account:fFCuZwED60:bindinfo"))
	keys, err := redis.Strings(con.Do("keys", "account:*:bindinfo"))
	if err != nil {
		fmt.Println("redis get error: ", err)
	}

	for _, key := range keys {
		value, err := redis.Values(con.Do("HGETALL", key))
		if err != nil {
			fmt.Println("redis get error: ", err)
			return
		}

		bindinfo := &AccountBindInfo{}
		if err := redis.ScanStruct(value, bindinfo); err != nil {
			fmt.Println("redis ScanStruct error: ", err)
			return
		}

		strs := strings.Split(key, ":")
		newkey := fmt.Sprintf("%v:%v:%v", strs[0], strs[2], strs[1])

		args := redis.Args{}.Add(newkey).AddFlat(bindinfo)
		_, err = con.Do("HMSET", args...)
		if err != nil {
			fmt.Println("hmset error: ", err)
			return
		}
	}
}

type AccountBasicInfo struct {
	AccountID  int32  `redis:"accountid,omitempty"`
	Token      string `redis:"token,omitempty"`
	OpenID     string `redis:"openid,omitempty"`
	SDKToken   string `redis:"sdktoken,omitempty"`
	PlatID     int32  `redis:"platid,omitempty"`
	DeviceInfo []byte `redis:"deviceinfo,omitempty"`
}

func TestCopyAccountBasicInfo(t *testing.T) {
	con := GetPool().Get()
	defer con.Close()

	//keys, err := redis.Strings(con.Do("keys", "account:100004023:basicinfo"))
	keys, err := redis.Strings(con.Do("keys", "account:*:basicinfo"))
	if err != nil {
		fmt.Println("redis get error: ", err)
	}

	for _, key := range keys {
		value, err := redis.Values(con.Do("HGETALL", key))
		if err != nil {
			fmt.Println("redis get error: ", err)
			return
		}

		info := &AccountBasicInfo{}
		if err := redis.ScanStruct(value, info); err != nil {
			fmt.Println("redis ScanStruct error: ", err)
			return
		}

		strs := strings.Split(key, ":")
		newkey := fmt.Sprintf("%v:%v:%v", strs[0], strs[2], strs[1])

		fmt.Println("oldkey", key, " ,newkey:", newkey)
		args := redis.Args{}.Add(newkey).AddFlat(info)
		_, err = con.Do("HMSET", args...)
		if err != nil {
			fmt.Println("hmset error: ", err)
			return
		}
	}
}
