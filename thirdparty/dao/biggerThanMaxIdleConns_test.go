package dao

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gotomicro/ego-component/egorm"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

// export EGO_DEBUG=true && go run main.go

type Test struct {
	Uid      int    `gorm:"not null" json:"uid"`
	Nickname string `gorm:"not null" json:"name"`
}

func (Test) TableName() string {
	return "test"
}

var config = `[logger.default]
enableAsync=false
[mysql.test]
#debug = true # ego重写gorm debug，打开后可以看到，配置名、地址、耗时、请求数据、响应数据
dsn = "root:nikki@(192.168.0.19:3306)/test?charset=utf8&parseTime=True&loc=Local&readTimeout=1s&timeout=1s&writeTimeout=3s"
onFail = "panic" # 失败后，直接fail fast，panic操作
connMaxLifetime = "300s"
maxIdleConns = 1
maxOpenConns = 100`

func TestMaxIdel(t *testing.T) {
	openDB()
	wg := sync.WaitGroup{}

	wg.Add(11)
	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				runDB()
			}
			wg.Done()
		}()
	}

	go func() {
		for {
			elog.Info("mysql info",
				elog.Int("MaxOpenConnections", gormDB.DB().Stats().MaxOpenConnections), //	db.maxOpen
				elog.Int("OpenConnections", gormDB.DB().Stats().OpenConnections),       //  db.numOpen
				elog.Int("InUse", gormDB.DB().Stats().InUse),                           // db.numOpen - len(db.freeConn)   len(db.freeConn)
				elog.Int64("MaxIdleClosed", gormDB.DB().Stats().MaxIdleClosed),         // gormDB.DB().Stats().MaxIdleClosed
			)
			if gormDB.DB().Stats().InUse == 0 && gormDB.DB().Stats().MaxIdleClosed > 0 {
				break
			}
			time.Sleep(100 * time.Microsecond)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("end")
}

var gormDB *egorm.Component

func openDB() error {
	econf.LoadFromReader(strings.NewReader(config), toml.Unmarshal)
	gormDB = egorm.Load("mysql.test").Build()
	return nil
}

func runDB() error {
	var test Test
	err := gormDB.Where("uid = ?", 2).Find(&test).Error
	fmt.Println("test ", test.Nickname, "err ", err)
	return err
}
