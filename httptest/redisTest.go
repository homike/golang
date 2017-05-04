package httptest

import (
	"fmt"
	"time"

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
			if _, err := c.Do("AUTH", "123456"); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	fmt.Println("Redis Init Success")
}

func GetPool() *redis.Pool {
	return RedisPool
}

const (
	key = "aaa"
)

func subscribe() {
	con := RedisPool.Get()
	defer con.Close()

	psc := redis.PubSubConn{con}
	psc.Subscribe("redChatRoom")

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s \n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println(v)
			return
		}
	}
}

func RunRedis() {
	c := RedisPool.Get()
	defer c.Close()

	var p1, p2 struct {
		Title  string `redis:"t"`
		Author string `redis:"a"`
		Body   string `redis:"b"`
	}

	p1.Title = "Example"
	p1.Author = "Gary"
	p1.Body = "Hello"

	if _, err := c.Do("HMSET", redis.Args{}.Add("id1").AddFlat(&p1)...); err != nil {
		fmt.Println(err)
		return
	}

	m := map[string]string{
		"title":  "Example2",
		"author": "Steve",
		"body":   "Map",
	}

	if _, err := c.Do("HMSET", redis.Args{}.Add("id2").AddFlat(m)...); err != nil {
		fmt.Println(err)
		return
	}

	for _, id := range []string{"id1", "id2"} {

		v, err := redis.Values(c.Do("HGETALL", id))
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := redis.ScanStruct(v, &p2); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%+v\n", p2)
	}

}

func InputData(i int) {
	c := RedisPool.Get()
	defer c.Close()

	name := fmt.Sprintf("czx%v", i)
	//fmt.Println(name)
	if _, err := c.Do("ZADD", "zsettest", i, name); err != nil {
		fmt.Println(err)
		return
	}
}

func GenData() {
	for i := 0; i < 100000; i++ {
		InputData(i)
	}
}

func TestRedis() {
	c := RedisPool.Get()
	defer c.Close()

	var countSec int64 = 0
	for i := 0; i < 1000; i++ {
		tBegin := time.Now().UnixNano()

		// rankList, err := redis.Values(c.Do("ZREVRANGEBYSCORE", "zsettest", 10, "-inf", "WITHSCORES", "LIMIT", 0, 10))
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// _ = rankList

		rankList, err := redis.Values(c.Do("ZRANGE", "zsettest", 10, 20, "WITHSCORES"))
		if err != nil {
			fmt.Println(err)
			return
		}
		_ = rankList

		tEnd := time.Now().UnixNano()

		countSec = countSec + (tEnd - tBegin)
	}

	fmt.Println("1000 : % us", countSec/1000)

}
