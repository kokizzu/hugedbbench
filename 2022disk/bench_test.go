package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/id64"
)

const DbName = `CockroachDB`
const GoRoutineCount = 100
const RecordsPerGoroutine = 400 // 100 conn x 400 insert/update
const SelectRepeat = 10         // 100 conn x 400 x 10 select
const SelectManyRepeat = 300    // 100 conn x 500 select 400 rows

// cockroach start-single-node --insecure

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

	conn, err := pgxpool.New(bg, pgUrl)
	L.PanicIf(err, `cannot connect db`)
	defer conn.Close()

	_, err = conn.Exec(bg, `CREATE TABLE IF NOT EXISTS bar1(id BIGINT PRIMARY KEY, foo VARCHAR(10) UNIQUE)`)
	L.PanicIf(err, `failed create table bar1`)

	_, err = conn.Exec(bg, `TRUNCATE table bar1`)
	L.PanicIf(err, `failed truncate table bar1`)

	wg := sync.WaitGroup{}

	t.Run(`insert`, func(t *testing.T) {
		for z := 0; z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z int) {
				for y := 0; y < RecordsPerGoroutine; y++ {
					_, err = conn.Exec(bg, `INSERT INTO bar1(id,foo) VALUES($1,$2)`, z*RecordsPerGoroutine+y, id64.SID())
					L.PanicIf(err, `failed insert to bar1`)
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	dur := time.Since(start).Seconds()
	fmt.Printf(DbName+` InsertOne: %.1f sec | %.1f rec/s`+"\n", dur, float64(GoRoutineCount*RecordsPerGoroutine)/dur)
	start = time.Now()

	row := conn.QueryRow(bg, `SELECT COUNT(1) FROM bar1`)
	count := 0
	err = row.Scan(&count)
	L.PanicIf(err, `failed query count/scan`)
	assert.Equal(t, GoRoutineCount*RecordsPerGoroutine, count)

	fmt.Println(DbName+` Count`, time.Since(start), count)
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

	dur = time.Since(start).Seconds()
	fmt.Printf(DbName+` UpdateOne: %.1f sec | %.1f rec/s`+"\n", dur, float64(GoRoutineCount*RecordsPerGoroutine)/dur)
	start = time.Now()

	row = conn.QueryRow(bg, `SELECT COUNT(1) FROM bar1`)
	count = 0
	err = row.Scan(&count)
	L.PanicIf(err, `failed query count/scan`)
	assert.Equal(t, GoRoutineCount*RecordsPerGoroutine, count)

	fmt.Println(DbName+` Count`, time.Since(start), count)
	start = time.Now()

	t.Run(`select`, func(t *testing.T) {
		for z := 0; z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z int) {
				var str string
				for y := 0; y < RecordsPerGoroutine; y++ {
					for x := 0; x < SelectRepeat; x++ {
						row := conn.QueryRow(bg, `SELECT foo FROM bar1 WHERE id=$1`, z*RecordsPerGoroutine+y)
						err := row.Scan(&str)
						L.PanicIf(err, `failed select bar1`)
					}
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	dur = time.Since(start).Seconds()
	fmt.Printf(DbName+` SelectOne: %.1f sec | %.1f rec/s`+"\n", dur, (GoRoutineCount*RecordsPerGoroutine)/dur)
	start = time.Now()

	totalScan := int64(0)
	totalQuery := int64(0)
	t.Run(`select-many`, func(t *testing.T) {
		query := fmt.Sprintf(`SELECT foo FROM bar1 WHERE id>=$1 AND id<$2 ORDER BY id LIMIT %d`, RecordsPerGoroutine)
		for z := 0; z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z int) {
				for x := 0; x < SelectManyRepeat; x++ {
					func() {
						var str string
						rows, err := conn.Query(bg, query, z*RecordsPerGoroutine, (z+1)*RecordsPerGoroutine)
						var res []string
						if err == nil {
							defer rows.Close()
							for rows.Next() {
								err := rows.Scan(&str)
								L.PanicIf(err, `failed scan bar1`)
								res = append(res, str)
							}
						}
						atomic.AddInt64(&totalScan, int64(len(res)))
						atomic.AddInt64(&totalQuery, 1)
						L.PanicIf(err, `failed select-many bar1`)
					}()
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	dur = time.Since(start).Seconds()
	fmt.Printf(DbName+` Select%d: %.1f sec | %.1f rec/s | %.1f queries/s`+"\n", RecordsPerGoroutine, dur, float64(totalScan)/dur, float64(totalQuery)/dur)
	start = time.Now()
}
