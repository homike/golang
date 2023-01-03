package thirdparty

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestRedisLuascript(t *testing.T) {
	RedisPool := &redis.Pool{
		MaxIdle:     1,
		IdleTimeout: 1,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	str := fmt.Sprintf(`local key = KEYS[1]; local ret = redis.call("ZRANGE",key,0,%v, "WITHSCORES");return ret`, -1)
	//str := fmt.Sprintf(`local ret = redis.call("ZRANGE", KEYS[1],0,%v, "WITHSCORES");return ret`, -1)
	//str := fmt.Sprintf(`local key = KEYS[1];return key`)

	script := redis.NewScript(1, str)

	fmt.Println("Redis Init Success")
	con := RedisPool.Get()

	//ids, err := redis.Int64s(script.Do(con, "userdata:7:test_dbmgr:dirtylist"))
	ret, err := script.Do(con, "10010001")
	if err != nil {
		fmt.Println("do script error:", err)
	}
	fmt.Println("[ids]", ret)
}
