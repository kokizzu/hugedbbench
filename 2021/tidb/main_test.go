package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/id64"
)

const DbName = `TiDB`
const GoRoutineCount = 1000
const RecordsPerGoroutine = 100

func TestDb(t *testing.T) {
	start := time.Now()
	defer func(start time.Time) {
		fmt.Println(DbName+` Total`, time.Since(start))
	}(start)

	bg := context.Background() // shared context
	myUrl := "%s:%s@tcp(%s:%d)/%s"
	myUrl = fmt.Sprintf(myUrl,
		`root`,
		``, // empty password
		`127.0.0.1`,
		4000,
		`test`,
	)

	conn, err := sql.Open("mysql", myUrl)
	L.PanicIf(err, `cannot connect db`)
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(1024)
	conn.SetMaxIdleConns(1024)
	defer conn.Close()

	_, err = conn.ExecContext(bg, `CREATE TABLE IF NOT EXISTS bar1(id BIGINT PRIMARY KEY, foo VARCHAR(10) UNIQUE)`)
	L.PanicIf(err, `failed create table bar1`)

	_, err = conn.ExecContext(bg, `TRUNCATE TABLE bar1`)
	L.PanicIf(err, `failed truncate table bar1`)

	wg := sync.WaitGroup{}

	t.Run(`insert`, func(t *testing.T) {
		for z := 0; z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z int) {
				for y := 0; y < RecordsPerGoroutine; y++ {
					uniq := id64.SID()
					_, err = conn.ExecContext(bg, `INSERT INTO bar1(id,foo) VALUES(?,?)`, I.ToStr(z*RecordsPerGoroutine+y), uniq)
					L.PanicIf(err, `failed insert to bar1`)
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	fmt.Println(DbName+` InsertOne`, time.Since(start))
	start = time.Now()

	row := conn.QueryRowContext(bg, `SELECT COUNT(1) FROM bar1`)
	count := 0
	err = row.Scan(&count)
	L.PanicIf(err, `failed query count/scan`)
	assert.Equal(t, GoRoutineCount*RecordsPerGoroutine, count)

	fmt.Println(DbName+` Count`, time.Since(start))
	start = time.Now()

	t.Run(`update`, func(t *testing.T) {
		for z := 0; z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z int) {
				for y := 0; y < RecordsPerGoroutine; y++ {
					uniq := id64.SID()
					_, err = conn.ExecContext(bg, `UPDATE bar1 SET foo=? WHERE id=?`, uniq, z*RecordsPerGoroutine+y)
					L.PanicIf(err, `failed update bar1`)
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	fmt.Println(DbName+` UpdateOne`, time.Since(start))
	start = time.Now()

	row = conn.QueryRowContext(bg, `SELECT COUNT(1) FROM bar1`)
	count = 0
	err = row.Scan(&count)
	L.PanicIf(err, `failed query count/scan`)
	assert.Equal(t, GoRoutineCount*RecordsPerGoroutine, count)

	fmt.Println(DbName+` Count`, time.Since(start))
	start = time.Now()

	t.Run(`select`, func(t *testing.T) {
		for z := 0; z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z int) {
				var str string
				for y := 0; y < RecordsPerGoroutine; y++ {
					row := conn.QueryRowContext(bg, `SELECT foo FROM bar1 WHERE id=?`, z*RecordsPerGoroutine+y)
					err := row.Scan(&str)
					L.PanicIf(err, `failed select bar1`)
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	fmt.Println(DbName+` SelectOne`, time.Since(start))
	start = time.Now()
}
