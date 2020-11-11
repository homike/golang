package methodhandler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Req struct {
	Req1 int
}

type Reply struct {
	Ret1 int
}

func Handler(req *Req) (*Reply, error) {
	return &Reply{
		Ret1: req.Req1 + 1,
	}, nil
}

func TestCallFunc(t *testing.T) {
	RegisterCall(METHOD_TEST, Handler)

	// 测试请求参数为nil
	{
		var req *Req
		_, err := ExecFunc(METHOD_TEST, []interface{}{interface{}(req)})

		assert.NotEqual(t, nil, err)
	}

	// 测试请求参数个数
	{
		input := 1
		req1 := interface{}(&Req{Req1: input})
		req2 := interface{}(&Req{Req1: input})
		_, err := ExecFunc(METHOD_TEST, []interface{}{req1, req2})

		assert.NotEqual(t, nil, err)
	}

	// 测试返回值
	{
		input := 1
		req := &Req{Req1: input}
		ret, err := ExecFunc(METHOD_TEST, []interface{}{req})
		assert.Equal(t, nil, err)

		reply, ok := ret.(*Reply)

		assert.Equal(t, true, ok)
		assert.Equal(t, input+1, reply.Ret1)
	}
}
