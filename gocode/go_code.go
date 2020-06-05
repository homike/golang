// 参考: https://github.com/golang/go/wiki/CodeReviewComments/5a40ba36d388ff1b8b2dd4c1c3fe820b8313152f
// 参考: https://golang.org/doc/effective_go.html

package main

import (
	"fmt"
	syslog "log"

	/*
		除了common模块, import 其他模块尽量不用 .
		import . 会导致程序易读性变差
	*/
	. "x5/common/errors"
	/*
		自动import的文件, 有时候可能是并不是预想中的, 注意检查

		### 项目中出现过, 想引用 game/util 的函数，自动引用了login/util, 出现了BUG的情况 ###
	*/)

/* 命名：遵循驼峰规则, 驼峰以单词为边界*/

/////////////////////////////////////////// Define ///////////////////////////////////////
/* 枚举定义  */
type ScoreLevel float64

const (
	ScoreLevel_Unknown = itoa
	ScoreLevel_Perfect
	ScoreLevel_Normal
)

/*
	常量定义, 以k 开头标注
	* 每个常量必须添加备注信息
	* 程序代码中不要出现数字相关的代码, 例如 if a > 1 {}
*/
const kPerfectScore float64 = 80.0 // 完美数值
const kNormalScore float64 = 60.0  // 正常数值

/*
	Interface 定义
	* interface的命名应该以er结尾，可以用来区别类型，更加容易分辨
*/
type Objecter interface {
	Name() string
	ID() float64
	IsGood() bool
}

/*
	日志级别
	* 每个操作的入口地址应该输出Info日志, 其他位置尽量控制 Info 日志的输出
	* 所有的日志是由错误级别的区分的，线上会开启Info级别
	* 凡事错误的日志，都需要携带玩家uid, 保证可能精准定位到玩家
	* 日志的输出尽量不用输出整个结构体，当结构体中有map时，并发读写可能会导致panic
	* 日志的打印要注意 打印对象为 nil 的情况, 尤其是Error日志, 可能到异常情况才能运行到

###  项目中出现过, Error 日志输出 err.Error(), 但是err为 nil, 导致程序崩溃的情况, 尽量不要使用err.Error()的打印  ###

###  项目中出现过, 由于输出日志过多，导致程序阻塞的情况  ###
*/
func (this *Upvalue) WriteComment(data []byte) bool {
	syslog.Printf("[Name: %v, ID: %v] len: %v, error: %v \n", this.name_, this.id_, data_len, err)
	return true
}

/////////////////////////////////////////// Struct ///////////////////////////////////////
/*
	Struct 定义
	* 成员变量私有化
	* 成员变量, 通过Get/Set来控制
	* 当不明确变脸类型时，应选更大取值范围的类型进行编写发，防止溢出
*/
type Upvalue struct {
	name string
}

/*
	Struct 方法
	* 针对自己的方法，使用self 来定义struct自身
	* Get 前缀不进行增加，更便于编写
	* Set 前缀增加，意义更加明确，能够更加正确的调用
*/
func (self *Upvalue) Name() string {
	return self.name
}
func (self *Upvalue) SetName(name string) bool {
	self.name = name
	return true
}

/*
	Struct 内部要保证struct数据的完整性
	* 外部接收内容， 不会影响内部逻辑
	* 外部不能修改内部内容， 只能通过回调来调用内部逻辑，并记录拷贝值

###  项目中出现过, 传递引用，导致配置数据被修改, 全服玩家的配置全部异常  ###
*/
func (this *Upvalue) TagListSlice([]string, bool) {
	ret_slice := []string{}
	for _, tag_name := range this.taglist_ {
		ret_slice = append(ret_slice, tag_name)

	}

	return ret_slice, len(ret_slice) > 0
}
func (this *Upvalue) ForeachTagListWithCallBackAllTags(cb func(int32, string) bool) bool {
	var ret bool = true
	for tag_id, tage_name := range this.taglist {
		ret = ret && cb(tag_id, tag_name)
	}
	return true
}

/////////////////////////////////////////// Function ///////////////////////////////////////
/*
	函数名
	* 函数名可以非常长，但是要保证意义正确
	* 异步函数需要在接口名, 上直接使用 Async 来进行标注

###  项目中出现过, 函数名为 GetUser(), 实际在user不存在时, 会createuser。新成员使用错误，出现BUG  ###
*/
func AsyncPrintObjecter(id int64, obj Objecter) {
	if obj == nil {
		syslog.Printf("print object is nil, id: %v", id)
		return
	}
	go func() {
		fmt.Printf("Objecter: [Name: %v, ID: %v]\n", obj.ID(), obj.Name())
	}()
}

/*
	函数返回值
	* 尽量采用值传递而不是引用传递
	* 返回的内容应该保证class / struct 内容安全, 所以应该赋值返回或者 new 返回
*/
func (this *Upvalue) UpvalueInfo() UpvalueIDName {
	return UpvalueIDName{
		name_: this.name_,
		id_:   this.id_,
	}
}

/* 当返回值为bool时, 函数名以Is/Has开头 */
func (this *Upvalue) IsGood() bool {
	return false
}

/*
	Error
	* 函数返回值为tuple返回，需要返回error内容，协助定位问题
	* 返回error后，外层需要处理，不能 ”下划线“ 无视

### 项目中出现过, 由于使用”下划线“无视error, 导致redis异常未能正确处理，玩家数据回档的问题 ###
*/
func (this *Upvalue) writeComment(data []byte) (num int, err error) {
	data = append(data, this.comment_...)
	data_len := len(data)

	/*
		else 分支
		减少使用else 的情况，固定逻辑分支判断，能够保证唯一性，且代码更加简单，保证逻辑可控
		不应该 使用
		if err != nil {
				// error handling
		} else {
					// normal code
		}
	*/
	if err != nil {
		// error handling
	}
	// normal code

	return data_len, nil
}

/////////////////////////////////////////// Other ///////////////////////////////////////
/*
	Slice
	* slice 在append 后, 优化slice自动扩容，此时的cap 和 len 已经不相等了, 需要关注

TODO: 待回顾
### 项目中出现过, 如下初始化，导致slice 的 cap 和预期的不一致，导致json.unmarshal 出现BUG  ###
*/
func GetSlice() []string {
	reward := LoginReward{0, 0}
	for i := 0; i < 40; i++ {
		user.Rewards = append(user.Rewards, reward)
	}
}

/*
	Float 内容的比较
	* 两数相减，按照范围取值判断
	* 直接比较两数 <, > 不要使用===, 以免误差太大

	当需要给一个函数内容添加备注时，可以考虑将备注信息抽离为一个函数进行编写，以免误差太大
*/
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
