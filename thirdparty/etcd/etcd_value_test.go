package etcd

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TablesServerinfoChannel_ChannelItem struct {
	Channel           string `json:"channel"`
	Name              string `json:"name"`
	LoginUrl          string `json:"login_url"`
	AndroidCdn        string `json:"android_cdn"`
	IosCdn            string `json:"ios_cdn"`
	AndroidApkVersion string `json:"android_apk_version"`
	IosApkVersion     string `json:"ios_apk_version"`
	AndroidApkUrl     string `json:"android_apk_url"`
	IosApkUrl         string `json:"ios_apk_url"`
	Timezone          int32  `json:"timezone"`
	Audit             int32  `json:"audit"`
}

func TestEtcdConfig(t *testing.T) {
	configs, err := GetServerConfig("192.168.0.18:2379,192.168.0.18:2381,192.168.0.18:2383", "x5config")
	if err != nil {
		t.Fatalf("GetServerConfig() error: %v", err)
	}

	for _, v := range configs {
		//fmt.Printf("[data]: file: %v, value: %v \n", v.Key, string(v.Value))

		if v.Key == "x5config/Tables.Serverinfo.Channel.Channel.json" {
			dataList := make([]*TablesServerinfoChannel_ChannelItem, 0)
			err = json.Unmarshal(v.Value, &dataList)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			fmt.Printf("Channel data: AndroidApkVersion: %v", dataList[0].AndroidApkVersion)
		}
	}
}
