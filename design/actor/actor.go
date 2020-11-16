package actor

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

var ErrSendFailed = errors.New("send error, client is shut down")
var ErrShutdown = errors.New("client is shut down")
var ErrDoFaield = errors.New("do something failed")

//----------------------------------- CALL -----------------------------------------
type Call struct {
	MsgKey  string
	MsgBody string
	Reply   interface{} // The reply from the function (*struct).
	Error   error       // After completion, the error status.
	Done    chan *Call  // Strobes when call is complete.
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		log.Println("discarding Call reply due to insufficient Done chan capacity")
	}
}

//----------------------------------- Client -----------------------------------------
type Client struct {
	requestList chan *Call
	SigClose    chan struct{}
	shutdown    bool // user has called Close
	m_CallMap   map[string]*CallFunc
}

func (client *Client) send(call *Call) {
	if client.shutdown {
		call.Error = ErrSendFailed
		call.done()
		return
	}

	fmt.Println("send ", call.MsgKey)
	select {
	case client.requestList <- call:
	default:
		fmt.Println("send messages error")
	}
}

func (client *Client) input() {
	var err error

loop:
	for err == nil {
		select {
		case call := <-client.requestList:
			// do something
			ret := client.doSomething(call.MsgKey)

			// 通过反射返回rely
			/*
				v := reflect.ValueOf(call.Reply)
				if v.Type().Kind() != reflect.Ptr {
					//dec.err = errors.New("gob: attempt to decode into a non-pointer")
					//return dec.err
				}
				if v.IsValid() {
					if v.Kind() == reflect.Ptr && !v.IsNil() {
						// That's okay, we'll store through the pointer.
					} else if !v.CanSet() {
						//return errors.New("gob: DecodeValue of unassignable value")
					}
				}
			*/

			call.Reply = &ret
			call.done()

		case <-client.SigClose:
			break loop
		default:
		}
	}

	// 通过shutdown来阻止
	fmt.Println("exit loop")
	fmt.Println(time.Now().Format("15:04:05.000"), " ,request len1: ", len(client.requestList))
	time.Sleep(5 * time.Second)

	client.shutdown = true

	fmt.Println(time.Now().Format("15:04:05.000"), " ,request len2: ", len(client.requestList))
	time.Sleep(1 * time.Second)

	fmt.Println("request len3: ", len(client.requestList))
	// 不能用for循环来关闭, 会阻塞等待
	select {
	case call := <-client.requestList:
		call.Error = ErrShutdown
		call.done()
	default:
	}
	//for call := range client.requestList {
	//}
}

func (client *Client) doSomething(index string) int {
	//fmt.Println("doSomething ", index)
	//time.Sleep(100 * time.Millisecond)
	return 1
}

func NewClient(uid int32) *Client {
	_ = uid
	client := &Client{
		requestList: make(chan *Call, 20),
		SigClose:    make(chan struct{}),
		shutdown:    false,
		m_CallMap:   make(map[string]*CallFunc),
	}
	go client.input()
	return client
}

type MsgHandle func(req interface{}) interface{}

var handles map[int]MsgHandle

func (client *Client) Register(cmd int, handle MsgHandle) {
	handles[cmd] = handle
}

func (client *Client) Handle(cmd int, req interface{}) interface{} {
	return handles[cmd](req)
}

func (client *Client) Close() error {
	if client.shutdown {
		return ErrShutdown
	}
	close(client.SigClose)
	return nil
}

// 异步调用
func (client *Client) Go(msgKey, msgBody string, reply interface{}) *Call {
	call := new(Call)
	call.MsgKey = msgKey
	call.MsgBody = msgBody
	call.Reply = reply
	call.Done = make(chan *Call, 1)

	client.send(call)
	return call
}

// 同步调用
func (client *Client) Call(msgKey, msgBody string, reply interface{}) error {
	call := <-client.Go(msgKey, msgBody, reply).Done
	return call.Error
}

type CallFunc struct {
	Func       interface{}
	FuncType   reflect.Type
	FuncVal    reflect.Value
	FuncParams string
}

func (this *Client) FindCall(funcName string) *CallFunc {
	funcName = strings.ToLower(funcName)
	fun, exist := this.m_CallMap[funcName]
	if exist == true {
		return fun
	}
	return nil
}

func (this *Client) RegisterCall(funcName string, call interface{}) {
	funcName = strings.ToLower(funcName)
	if this.FindCall(funcName) != nil {
		fmt.Println(fmt.Errorf("cannot find call"))
	}

	this.m_CallMap[funcName] = &CallFunc{Func: call, FuncVal: reflect.ValueOf(call), FuncType: reflect.TypeOf(call), FuncParams: reflect.TypeOf(call).String()}
}

func (this *Client) ExecFunc(funcName string, param interface{}) []reflect.Value {
	params := make([]interface{}, 0)
	params = append(params, param)

	pFunc := this.FindCall(funcName)
	if pFunc != nil {
		f := pFunc.FuncVal
		k := pFunc.FuncType
		strParams := pFunc.FuncParams

		if k.NumIn() != len(params) {
			log.Printf("func [%s] can not call, func params [%s], params [%v]", funcName, strParams, params)
			return nil
		}

		if len(params) >= 1 {
			in := make([]reflect.Value, len(params))
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
			}

			ret := f.Call(in)
			return ret
		} else {
			log.Printf("func [%s] params at least one context", funcName)
			//f.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}
	}

	return nil
}
