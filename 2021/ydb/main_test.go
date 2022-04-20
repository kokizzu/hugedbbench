package main

import (
	"context"
	"fmt"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/id64"
	"github.com/stretchr/testify/assert"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
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

	db, err := ydb.New(
		connectCtx,
		ydb.WithConnectionString("grpc://127.0.0.1:2136/?database=/local"),
		ydb.WithAnonymousCredentials(),
	)
	L.PanicIf(err, `connect.New`)
	ctx := context.Background()
	defer db.Close(ctx)

	session, err := db.Table().CreateSession(ctx)
	L.PanicIf(err, `db.Table().CreateSession`)
	defer session.Close(ctx)

	txc := table.TxControl(
		table.BeginTx(table.WithSerializableReadWrite()),
		table.CommitTx(),
	)

	// start benchmarking

	err = db.Table().Do(
		ctx, func(ctx context.Context, s table.Session) error {
			tableName := path.Join(db.Name(), `bar1`)
			err := s.DropTable(ctx, tableName)
			L.Print(err) // ignore error if first time
			return s.CreateTable(ctx, tableName,
				options.WithColumn(`id`, types.Optional(types.TypeUint64)),
				options.WithColumn(`foo`, types.Optional(types.TypeUTF8)),
				options.WithPrimaryKeyColumn(`id`),
				options.WithIndex(`idx_foo`, options.WithIndexColumns(`foo`)),
			)
		})
	L.PanicIf(err, `failed create table bar1`) // no UNIQUE index supported

	wg := sync.WaitGroup{}

	// https://ydb.tech/en/docs/reference/ydb-sdk/example/go/

	t.Run(`insert`, func(t *testing.T) {
		for z := uint64(0); z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z uint64) {
				session, err := db.Table().CreateSession(ctx)
				L.PanicIf(err, `db.Table().CreateSession`)
				defer session.Close(ctx)
				for y := uint64(0); y < RecordsPerGoroutine; y++ {
					uniq := id64.SID()
					_, _, err = session.Execute(ctx, txc, `
Declare $id AS Uint64;
Declare $foo AS Utf8;
INSERT INTO bar1(id,foo) VALUES($id,$foo)`,
						table.NewQueryParameters(
							table.ValueParam(`$id`, types.Uint64Value(z*RecordsPerGoroutine+y)),
							table.ValueParam(`$foo`, types.UTF8Value(uniq)),
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

	count := func() {
		_, res, err := session.Execute(ctx, txc, `SELECT COUNT(1) AS cou FROM bar1`, nil)
		L.PanicIf(err, `failed select count(1) from bar1`)
		defer res.Close()
		count := uint64(0)
		err = res.NextResultSetErr(ctx)
		L.PanicIf(err, `failed next result set`)
		if res.NextRow() {
			err = res.ScanNamed(named.Required(`cou`, &count))
			L.PanicIf(err, `failed query count/scan`)
		}
		assert.Equal(t, GoRoutineCount*RecordsPerGoroutine, int(count))
	}
	count()

	fmt.Println(DbName+` Count`, time.Since(start))
	start = time.Now()

	t.Run(`update`, func(t *testing.T) {
		for z := uint64(0); z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z uint64) {
				session, err := db.Table().CreateSession(ctx)
				L.PanicIf(err, `db.Table().CreateSession`)
				defer session.Close(ctx)
				for y := uint64(0); y < RecordsPerGoroutine; y++ {
					uniq := id64.SID()
					_, _, err = session.Execute(ctx, txc, `
Declare $id AS Uint64;
Declare $foo AS Utf8;
UPDATE bar1 SET foo=$foo WHERE id=$id`,
						table.NewQueryParameters(
							table.ValueParam(`$id`, types.Uint64Value(z*RecordsPerGoroutine+y)),
							table.ValueParam(`$foo`, types.UTF8Value(uniq)),
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

	count()

	fmt.Println(DbName+` Count`, time.Since(start))
	start = time.Now()

	t.Run(`select`, func(t *testing.T) {
		for z := uint64(0); z < GoRoutineCount; z++ {
			wg.Add(1)
			go func(z uint64) {
				session, err := db.Table().CreateSession(ctx)
				L.PanicIf(err, `db.Table().CreateSession`)
				defer session.Close(ctx)
				var str *string // must add pointer
				for y := uint64(0); y < RecordsPerGoroutine; y++ {
					_, res, err := session.Execute(ctx, txc, `
Declare $id AS Uint64;
SELECT foo FROM bar1 WHERE id=$id`,
						table.NewQueryParameters(
							table.ValueParam(`$id`, types.Uint64Value(z*RecordsPerGoroutine+y)),
						),
					)
					err = res.NextResultSetErr(ctx)
					L.PanicIf(err, `failed next result set`)
					if res.NextRow() {
						err = res.Scan(&str) // works when str is a pointer to string
						//err = res.Scan(&str) //  {s:"scan row failed: type *string is not optional! use double pointer or sql.Scanner."},
						//err = res.ScanNamed(named.Required(`foo`, &str)) // {s:"scan row failed: incorrect source types PRIMITIVE_TYPE_ID_UNSPECIFIED"},
						L.PanicIf(err, `failed scan bar1`)
					}
					_ = res.Close()
				}
				wg.Done()
			}(z)
		}
		wg.Wait()
	})

	fmt.Println(DbName+` SelectOne`, time.Since(start))
	start = time.Now()

}

/* benchmark result:
YDB InsertOne 1m2.203654718s
YDB Count 15.139201ms
YDB UpdateOne 1m11.049503884s
YDB Count 25.023836ms
YDB SelectOne 11.580145895s
YDB Total 2m24.874077468s
6.4G    2021/ydb/ydb_data/
*/
