package basetest

/*****************************************************************************************
* 测试1， 直接存储接口, 打印输出为空
*****************************************************************************************/

// type TestStruct struct {
// 	TestData map[int]TestIf `json:"data"`
// }

// type TestIf interface {
// 	DoTest()
// }

// // Father
// type Father struct {
// 	Var1 uint `json:"var1"`
// 	Var2 uint `json:"var2"`
// }

// func NewFather(v1, v2 uint) *Father {
// 	return &Father{
// 		Var1: v1,
// 		Var2: v2,
// 	}
// }

// func (f *Father) DoTest() {
// 	fmt.Println("value: ", f.Var1)
// }

// // Child
// type Child struct {
// 	Var1 uint `json:"var1"`
// 	Var2 uint `json:"var2"`
// 	Var3 uint `json:"var3"`
// 	Var4 uint `json:"var4"`
// }

// func NewChild(v1, v2, v3, v4 uint) *Child {
// 	return &Child{
// 		Var1: v1,
// 		Var2: v2,
// 		Var3: v3,
// 		Var4: v4,
// 	}
// }

// func (c *Child) DoTest() {
// 	fmt.Println("value: ", c.Var1)
// }

// func RunIF() {
// 	var v1 TestIf = (interface{})(NewChild(1, 2, 3, 4)).(TestIf)
// 	var v2 TestIf = (interface{})(NewFather(1, 2)).(TestIf)

// 	tMap := make(map[int]TestIf)
// 	tMap[1] = v1
// 	tMap[2] = v2

// 	data1 := TestStruct{
// 		TestData: tMap,
// 	}

// 	bData, err := json.Marshal(data1)
// 	if err != nil {
// 		return
// 	}

// 	fmt.Println("data: ", string(bData))
// 	// data:  {"data":{"1":{"var1":1,"var2":2,"var3":3,"var4":4}, "2":{"var1":1,"var2":2}}}
// }

/*****************************************************************************************
* 测试2， 直接存储接口, 打印输出为空
*****************************************************************************************/

type FatherIF interface {
	DoTest() uint
	DoTest2() uint
}

type Father struct {
	Var1 uint `json:"var1"`
	Var2 uint `json:"var2"`
}

func (mt *Father) DoTest() uint {
	return mt.Var1
}

func (mt *Father) DoTest2() uint {
	return mt.Var2
}

type Child struct {
	Var3    uint `json:"var3"`
	Var4    uint `json:"var4"`
	*Father `json:"father"`
}

func NewChild(v1, v2, v3, v4 uint) *Child {
	return &Child{
		Father: &Father{
			Var1: v1,
			Var2: v2,
		},
		Var3: v3,
		Var4: v4,
	}
}

func (mt *Child) DoTest3() uint {
	return mt.Var3
}

type TestStruct struct {
	TestData map[int]*FatherIF `json:"data"`
}

func RunIF() {

	// Test2
	// v1 := new(FatherIF)
	// *v1 = NewChild(1, 2, 3, 4)

	// ret1 := (*v1).DoTest()
	// ret2 := (*v1).DoTest2()

	// c, ok := (interface{})(*v1).(*Child)
	// if !ok {

	// }
	// ret3 := c.DoTest3()
	// fmt.Println("resut: ", ret1, ret2, ret3)

	// Test1
	// tMap := make(map[int]*TestIf)
	// tMap[1] = v1

	// data1 := TestStruct{
	// 	TestData: tMap,
	// }

	// bData, err := json.Marshal(data1)
	// if err != nil {
	// 	return
	// }

	// fmt.Println("data: ", string(bData))

	// {"data":{"1":{}}}
}
