package thirdparty

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
)

func Test_MapNil(t *testing.T) {
	test := &Test{Infos: make(map[int32]*TestInfo)}
	//test.Infos[1] = &TestInfo{InfoMap: make(map[int32]bool)}
	test.Infos[1] = nil
	pbTest, _ := proto.Marshal(test)

	unMarshalTest := &Test{}
	proto.Unmarshal(pbTest, unMarshalTest)

	fmt.Println(unMarshalTest)
	fmt.Println(unMarshalTest.Infos[1].InfoMap[1])

}
