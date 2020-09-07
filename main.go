package main

import (
	"fmt"
	"gotest/basetest"
)

// 关卡ID规则
// 10010101
// 1001     章节ID
// 100101   具体ID
// 10010101 难度

func GetLevelType(svrLevelID int32) int32 {
	return svrLevelID / 10000000
}

func GetChapter(svrLevelID int32) int32 {
	return svrLevelID / 10000
}

func GetCommonLevelID(svrLevelID int32) int32 {
	return svrLevelID / int32(100)
}

func GetDifficulty(svrLevelID int32) int32 {
	return svrLevelID % int32(100)
}

func ParseSvrLevelID(svrLevelID int32) (int32, int32, int32) {
	return GetLevelType(svrLevelID), GetCommonLevelID(svrLevelID), GetDifficulty(svrLevelID)
}

func GenSvrLevelID(levelID, diff int32) int32 {
	return levelID*10 + diff
}

func main() {
	basetest.RunReturn()

	ltype, lcommonID, ldiff := ParseSvrLevelID(10020304)
	fmt.Printf("parse levelID : %v, %v, %v \n", ltype, lcommonID, ldiff)
}
