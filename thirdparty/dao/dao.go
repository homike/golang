package dao

import (
	"database/sql"
	"fmt"
	"runtime"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var nuanDB *sql.DB
var SavingQueue = make(chan int, 10000)

func RunExecSql(id, age int, name string) {
	fmt.Println("exec start", id)
	SavingQueue <- id

	defer func() {
		<-SavingQueue
		fmt.Println("exec end", id, "len(queue)", len(SavingQueue))
	}()

	_, err := nuanDB.Exec(fmt.Sprintf("update user set name = '%s', age = %d WHERE id = %d", name, age, id))
	if err != nil {
		fmt.Println("exec error", err)
		return
	}
}

func RunInsertSql(id, age int, name string) {
	fmt.Println("exec start")
	_, err := nuanDB.Exec(fmt.Sprintf("INSERT INTO user(name, age, id) VALUES('%s', %d, %d)", name, age, id))
	if err != nil {
		fmt.Println("exec error", err)
		return
	}
	fmt.Println("exec end")
}

func RunDao() {
	//SavingQueue = make(chan uint, 10000)
	dbdns := "root:123456@(192.168.0.243:3306)/test?parseTime=true&loc=Local&charset=utf8"
	maxDBConn := 100
	var err error = nil
	if nuanDB, err = sql.Open("mysql", dbdns); err != nil {
		fmt.Printf("sql.Open(\"mysql\", %s) failed (%v)", dbdns, err)
		return
	}

	nuanDB.SetMaxIdleConns(maxDBConn / 2)
	nuanDB.SetMaxOpenConns(maxDBConn)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			fmt.Println(fmt.Sprintf("NumGoroutine:%v", runtime.NumGoroutine()))
		}
	}()

	fmt.Println("Dao Start1")
	var wgroup sync.WaitGroup
	for i := 1; i < 1000; i++ {
		wgroup.Add(1)
		go func(id, age int, name string) {
			RunExecSql(id, age, name)
			wgroup.Done()
		}(i, i, fmt.Sprintf("name%d", i))
	}
	wgroup.Wait()

	fmt.Println("Dao Start2")
	time.Sleep(10 * time.Second)
	for i := 1000; i < 2000; i++ {
		wgroup.Add(1)
		go func(id, age int, name string) {
			RunExecSql(id, age, name)
			wgroup.Done()
		}(i, i, fmt.Sprintf("name%d", i))
	}
	wgroup.Wait()

	fmt.Println("Dao Start3")
	time.Sleep(10 * time.Second)
	for i := 2000; i < 3000; i++ {
		wgroup.Add(1)
		go func(id, age int, name string) {
			RunExecSql(id, age, name)
			wgroup.Done()
		}(i, i, fmt.Sprintf("name%d", i))
	}
	wgroup.Wait()

	fmt.Println("Dao End----------------")

	for {
	}
}
