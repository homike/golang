// Createtime: 2020/01/01 [yyyy/mm/dd]
// Author: moyu [git commit name]
// Commmit: 文件介绍
package main

import (
	"fmt"
	syslog "log"
)

type ScoreLevel float64

const (
	ScoreLevel_Unknown = itoa
	ScoreLevel_Perfect
	ScoreLevel_Normal
)

const kPerfectScore float64 = 80.0 // 完美数值
const kNormalScore float64 = 60.0  // 正常数值

type Objecter interface {
	Name() string
	ID() float64
	IsGood() bool
}

type UpvalueIDName struct {
	name_ string
	id_   float64
}

type Upvalue struct {
	name_    string
	id_      float64
	taglist_ map[int32]string
	score_   float64
	comment_ []byte
}

func (this *Upvalue) Name() string {
	return this.name_
}

func (this *Upvalue) SetName(name string) bool {
	this.name_ = name
	return true
}

func (this *Upvalue) UpvalueInfo() UpvalueIDName {
	return UpvalueIDName{
		name_: this.name_,
		id_:   this.id_,
	}
}

func (this *Upvalue) IsGood() bool {
	scorelevel := this.GetScoreLevel()
	if scorelevel == ScoreLevel_Perfect {
		return true
	}
	return false
}

func (this *Upvalue) GetScoreLevel() ScoreLevel {
	ScoreLevel := ScoreLevel_Unknown
	switch {
	case this.score_ < kNormalScore:
		ret = ScoreLevel_Unknown
	case this.score_ >= kNormalScore && this.score_ < kPerfectScore:
		ret = ScoreLevel_Normal
	case this.score_ >= kPerfectScore:
		ret = ScoreLevel_Perfect
	default:
		syslog.Printf("Misunderstand logic")
	}
	return ret
}

func (this *Upvalue) TagListSlice([]string, bool) {
	ret_slice := []string{}
	for _, tag_name := range this.taglist_ {
		ret_slice = append(ret_slice, tag_name)

	}

	return ret_slice, len(ret_slice) > 0
}

func (this *Upvalue) ForeachTagListWithCallBackAllTags(cb func(int32, string) bool) bool {
	ret := true
	for tag_id, tag_name := range this.taglist_ {
		ret = ret & cb(tag_id, tag_name)
	}
	return ret
}

func (this *Upvalue) WriteComment(data []byte) bool {
	len, err := this.writeComment(data)
	if err != nil {
		syslog.Printf("[Name: %v, ID: %v] len: %v, error: %v \n", this.name_, this.id_, data_len, err)
		return false
	}
	return true
}

func (this *Upvalue) writeComment(data []byte) (num int, err error) {
	data = append(data, this.comment_...)
	data_len := len(data)
	if data_len <= 0 {
		return data_len, fmt.Errorf("data len error")
	}

	return data_len, nil
}

func AsyncPrintObjecter(id int64, obj Objecter) {
	if obj == nil {
		syslog.Printf("print object is nil, id: %v", id)
		return
	}
	go func() {
		fmt.Printf("Objecter: [Name: %v, ID: %v]\n", obj.ID(), obj.Name())
	}()
}
