package thirdparty

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"
)

type SimConfig struct {
	LoginAddr string `xml:"login_addr"`
	GameAddr  string `xml:"game_addr"`
	ClientVer string `xml:"client_ver"`
	BaseVer   string `xml:"base_ver"`
}

func readFile(filename string) (*SimConfig, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("ReadFile:", err.Error())
		return nil, err
	}

	v := SimConfig{}
	err = xml.Unmarshal(bytes, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}

	return &v, nil

}

func TestXML(t *testing.T) {
	SimCfg, err := readFile("sim.xml")

	if err != nil {
		fmt.Println("readFile:", err.Error())

		SimCfg = &SimConfig{
			LoginAddr: "https://192.168.0.140:8380",
			GameAddr:  "http://192.168.0.140:8680",
			ClientVer: "1.2.5",
			BaseVer:   "1.6.0",
		}
	}

	_ = SimCfg
}
