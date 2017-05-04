package interfacetest

import (
	"fmt"
)

type RobotFunc interface {
	InitCase(sceneId int)
}

type Robot struct {
	RobotIndex int
	RobotFunc
}

func FanInRobot(rFunc interface{}) chan *Robot {
	value, ok := rFunc.(RobotFunc)
	if !ok {
		fmt.Println("convert error: %v", value)
		return nil
	}
	fmt.Println("FanInRobot")

	robots := make(chan *Robot, 2000)
	for i := 0; i < 1; i++ {
		//name := RobotName //fmt.Sprintf("robot%v", i)
		go func() {
			fmt.Println("FanInRobot 1")
			robot := &Robot{
				RobotIndex: i,
				RobotFunc:  value,
			}
			value.InitCase(1)
			robot.RobotFunc.InitCase(2)

			robots <- robot
		}()
	}

	return robots
}

//===============================================
type MRobot struct {
	Robot

	Name string
}

func (r *MRobot) InitCase(sceneID int) {
	fmt.Println("InitCase", sceneID)
	r.Name = fmt.Sprintf("robot%v", sceneID)
	//fmt.Println(sceneID, r.Name)
}
