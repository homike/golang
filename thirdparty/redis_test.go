package thirdparty

import (
	"fmt"
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

type AccountBasicInfo struct {
	AccountID  int32  `redis:"accountid"`
	Permission int32  `redis:"permission" default:"1"`
	Token      string `redis:"token"`
	OpenID     string `redis:"open_id"`
}

// 测试redigo/redis 序列化数据
func TestRedisHgetAll(t *testing.T) {
	con := GetPool().Get()

	value, err := redis.Values(con.Do("HGETALL", "xxx"))
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
func TestRedisEmpty(t *testing.T) {
	con := GetPool().Get()

	_, err := redis.Values(con.Do("HGETALL", "xxx"))
	//value, err := redis.Values(con.Do("get", key))
	if err != redis.ErrNil {
		fmt.Println("redis get error: ErrNil")
	} else {
		fmt.Println("redis get error: ", err)
	}
}
