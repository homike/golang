package dao

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// export EGO_DEBUG=true && go run main.go
func TestMaxIdel2(t *testing.T) {
	openTestDB()
	wg := sync.WaitGroup{}

	wg.Add(101)
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
		go func() {
			for i := 0; i < 1000; i++ {
				runTestDB()
			}
			wg.Done()
		}()
	}

	go func() {
		for {
			fmt.Printf(
				"%v, mysql info, MaxOpenConnections: %v, OpenConnections: %v, InUse: %v, MaxIdleClosed: %v \n",
				time.Now().Format("15:04:05.000"),
				testDB.Stats().MaxOpenConnections, //	db.maxOpen
				testDB.Stats().OpenConnections,    //  db.numOpen
				testDB.Stats().InUse,              // db.numOpen - len(db.freeConn)   len(db.freeConn)
				testDB.Stats().MaxIdleClosed,      // testDB.Stats().MaxIdleClosed
			)
			if testDB.Stats().InUse == 0 && testDB.Stats().MaxIdleClosed > 0 {
				break
			}
			time.Sleep(100 * time.Microsecond)
		}
		wg.Done()
	}()
	wg.Wait()
}

var testDB *sql.DB

func openTestDB() error {
	dbdns := "root:nikki@(192.168.0.19:3306)/test?parseTime=true&loc=Local&charset=utf8"
	var err error = nil
	if testDB, err = sql.Open("mysql", dbdns); err != nil {
		fmt.Printf("sql.Open(\"mysql\", %s) failed (%v)", dbdns, err)
		return err
	}

	testDB.SetMaxOpenConns(100)
	testDB.SetMaxIdleConns(1)

	return nil
}

func runTestDB() error {
	rows, err := testDB.Query(fmt.Sprintf("select name from test  WHERE uid = %d", 2))
	if err != nil {
		fmt.Println("exec error", err)
		return err
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			fmt.Printf("rows.Scan error: %v \n", err)
			return err
		}
	}

	return nil
}
