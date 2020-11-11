package handler

import (
	"fmt"
	"testing"

	"gopkg.in/square/go-jose.v2/json"
)

type AReq struct {
	Uid  int64  `json:"Uid"`
	Data string `json:"data"`
}

type AService struct {
}

func (self *AService) Request(s *Session, msg *AReq) error {
	fmt.Println("Test1Handler.MessageHandler.Request, msg:  ", msg)
	return nil
}

type BReq struct {
	Uid  int64  `json:"Uid"`
	Data string `json:"data"`
}

type BService struct {
}

func (self *BService) Request(s *Session, msg *BReq) error {
	fmt.Println("Test2Handler.MessageHandler.Request, msg:  ", msg)
	return nil
}

func TestService(t *testing.T) {
	sA := &AService{}
	err := Register(sA)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	sB := &BService{}
	err = Register(sB)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	// A
	{
		req := AReq{
			Uid:  111,
			Data: "aaa",
		}
		reqBytes, _ := json.Marshal(req)
		ProcessMessage("AService.Request", reqBytes)
	}
	// B
	{
		req := BReq{
			Uid:  222,
			Data: "bbb",
		}
		reqBytes, _ := json.Marshal(req)
		ProcessMessage("BService.Request", reqBytes)
	}
}
