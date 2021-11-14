package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/id64"
)

const DbName = `CockroachDB`
const GoRoutineCount = 1000
const RecordsPerGoroutine = 100

// docker-compose -f docker-compose-single.yaml up --remove-orphans
// docker-compose -f docker-compose-multi.yaml up --remove-orphans

func TestDb(t *testing.T) {
	start := time.Now()
	defer func(start time.Time) {
		fmt.Println(DbName+` Total`, time.Since(start))
	}(start)
	
	bg := context.Background() // shared context
	pgUrl := "postgres://%s:%s@%s:%d/%s"
	pgUrl = fmt.Sprintf(pgUrl,
		`root`,
		``, // empty password
		`127.0.0.1`,
		26257,
		`defaultdb`,
	)

	conn, err := pgxpool.Connect(bg, pgUrl)
	L.PanicIf(err, `cannot connect db`)
	defer conn.Close()

	_, err = conn.Exec(bg, `CREATE TABLE IF NOT EXISTS bar1(id BIGINT PRIMARY KEY, foo VARCHAR(10) UNIQUE)`)
	L.PanicIf(err, `failed create table bar1`)

	_, err = conn.Exec(bg, `TRUNCATE TABLE bar1`)
	L.PanicIf(err, `failed truncate table bar1`)

	wg := sync.WaitGroup{}

	t.Run(`insert`, func(t *testing.T) {
		for z := 0; z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z int) {
				for y := 0; y < RecordsPerGoroutine; y++ {
					_, err = conn.Exec(bg, `INSERT INTO bar1(id,foo) VALUES($1,$2)`, I.ToStr(z*RecordsPerGoroutine+y), id64.SID())
					L.PanicIf(err, `failed insert to bar1`)
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})
	
	fmt.Println(DbName+` InsertOne`, time.Since(start))
	start = time.Now()

	row := conn.QueryRow(bg, `SELECT COUNT(1) FROM bar1`)
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
					_, err = conn.Exec(bg, `UPDATE bar1 SET foo=$1 WHERE id=$2`, id64.SID(), z*RecordsPerGoroutine+y)
					L.PanicIf(err, `failed update bar1`)
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})
	
	fmt.Println(DbName+` UpdateOne`, time.Since(start))
	start = time.Now()

	row = conn.QueryRow(bg, `SELECT COUNT(1) FROM bar1`)
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
					row := conn.QueryRow(bg, `SELECT foo FROM bar1 WHERE id=$1`, z*RecordsPerGoroutine+y)
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
