package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/id64"
	"github.com/stretchr/testify/assert"
	"github.com/yandex-cloud/ydb-go-sdk/v2"
	"github.com/yandex-cloud/ydb-go-sdk/v2/connect"
	"github.com/yandex-cloud/ydb-go-sdk/v2/table"
)

// docker-compose -f docker-compose-single.yaml up

const DbName = `YDB`
const GoRoutineCount = 1000
const RecordsPerGoroutine = 100

// curl https://storage.yandexcloud.net/yandexcloud-ydb/install.sh | bash
// ydb -e grpc://localhost:2136 -d /local scheme ls

func TestDb(t *testing.T) {
	start := time.Now()
	defer func(start time.Time) {
		fmt.Println(DbName+` Total`, time.Since(start))
	}(start)

	connectCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := connect.New(
		connectCtx,
		connect.MustConnectionString(
			"grpc://127.0.0.1:2136/?database=/local",
		),
	)
	L.PanicIf(err, `connect.New`)
	defer db.Close()

	ctx := context.Background()
	session, err := db.Table().CreateSession(ctx)
	L.PanicIf(err, `db.Table().CreateSession`)
	defer session.Close(ctx)

	txc := table.TxControl(
		table.BeginTx(table.WithSerializableReadWrite()),
		table.CommitTx(),
	)

	// start benchmarking

	_, _, _ = session.Execute(ctx, txc, `DROP TABLE IF EXISTS bar1`, nil)

	_, _, err = session.Execute(ctx, txc, `
CREATE TABLE bar1(
	id Uint64, 
	foo Utf8, 
	INDEX idx_foo GLOBAL ON (foo),
	PRIMARY KEY(id)
)`, nil)
	L.PanicIf(err, `failed create table bar1`) // no UNIQUE index supported

	wg := sync.WaitGroup{}

	t.Run(`insert`, func(t *testing.T) {
		for z := uint64(0); z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z uint64) {
				for y := uint64(0); y < RecordsPerGoroutine; y++ {
					uniq := id64.SID()
					_, _, err = session.Execute(ctx, txc, `INSERT INTO bar1(id,foo) VALUES($id,$foo)`,
						table.NewQueryParameters(
							table.ValueParam(`$id`, ydb.Uint64Value(z*RecordsPerGoroutine+y)),
							table.ValueParam(`$foo`, ydb.UTF8Value(uniq)),
						),
					)
					L.PanicIf(err, `failed insert to bar1`)
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	fmt.Println(DbName+` InsertOne`, time.Since(start))
	start = time.Now()

	_, row, err := session.Execute(ctx, txc, `SELECT COUNT(1) AS cou FROM bar1`, nil)
	L.PanicIf(err, `failed select count(1) from bar1`)
	count := 0
	if row.NextResultSet(ctx, `cou`) {
		err = row.Scan(&count)
		L.PanicIf(err, `failed query count/scan`)
	}
	L.PanicIf(row.Err(), `row.NextResultSet`)
	assert.Equal(t, GoRoutineCount*RecordsPerGoroutine, count)

	fmt.Println(DbName+` Count`, time.Since(start))
	start = time.Now()

	t.Run(`update`, func(t *testing.T) {
		for z := uint64(0); z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z uint64) {
				for y := uint64(0); y < RecordsPerGoroutine; y++ {
					uniq := id64.SID()
					_, _, err = session.Execute(ctx, txc, `UPDATE bar1 SET foo=$foo WHERE id=$id`,
						table.NewQueryParameters(
							table.ValueParam(`$id`, ydb.Uint64Value(z*RecordsPerGoroutine+y)),
							table.ValueParam(`$foo`, ydb.UTF8Value(uniq)),
						),
					)
					L.PanicIf(err, `failed update bar1`)
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	fmt.Println(DbName+` UpdateOne`, time.Since(start))
	start = time.Now()

	_, row, err = session.Execute(ctx, txc, `SELECT COUNT(1) AS cou FROM bar1`, nil)
	L.PanicIf(err, `failed select count(1) from bar1`)
	count = 0
	if row.NextResultSet(ctx, `cou`) {
		err = row.Scan(&count)
		L.PanicIf(err, `failed query count/scan`)
	}
	L.PanicIf(row.Err(), `row.NextResultSet`)
	assert.Equal(t, GoRoutineCount*RecordsPerGoroutine, count)

	fmt.Println(DbName+` Count`, time.Since(start))
	start = time.Now()

	t.Run(`select`, func(t *testing.T) {
		for z := uint64(0); z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z uint64) {
				var str string
				for y := uint64(0); y < RecordsPerGoroutine; y++ {
					_, row, err := session.Execute(ctx, txc, `SELECT foo FROM bar1 WHERE id=$id`,
						table.NewQueryParameters(
							table.ValueParam(`$id`, ydb.Uint64Value(z*RecordsPerGoroutine+y)),
						),
					)
					L.PanicIf(err, `failed select foo from bar1`)
					err = row.Scan(&str)
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
