package loopVsWhereIn

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/kamalshkeir/pgdriver"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"

	"loopVsWhereIn/mTest"
)

const total = 100000
const limit = 1000 // if changed, change also limit1k in tarantool_test.go
const cores = 32

// pgx

const pgxPostgresConnStr = `postgres://pgroot:password@localhost:5432`
const pgxCockroachConnStr = `postgres://root@localhost:26257`

type PgxTestTable struct {
	Id      uint64
	Content string
}

const pgxTableName = `test_table3`

const pgxMigrateSql = `
CREATE TABLE IF NOT EXISTS test_table3 (
	id BIGINT NOT NULL,
	content TEXT NOT NULL,
	CONSTRAINT test_table3_pkey PRIMARY KEY (id),
	UNIQUE(content)
);`

func TestMain(m *testing.M) {
	var err error

	log.Println(`postgres`)
	{
		log.Println(`pgx`)
		{
			pgxPostgres, err = pgxpool.New(context.Background(), pgxPostgresConnStr)
			L.PanicIf(err, `pgx.Connect`)
			_, err = pgxPostgres.Exec(context.Background(), pgxMigrateSql)
			L.PanicIf(err, `pgx.Exec`)
		}
	}

	log.Println(`cockroach`)
	{
		log.Println(`pgx`)
		{
			pgxCockroach, err = pgxpool.New(context.Background(), pgxCockroachConnStr)
			L.PanicIf(err, `pgx.Connect`)
			_, err = pgxCockroach.Exec(context.Background(), pgxMigrateSql)
			L.PanicIf(err, `pgx.Exec`)
		}
	}

	log.Println(`tarantool`)
	{
		taran = &Tt.Adapter{Connection: mTest.ConnectTarantool(), Reconnect: mTest.ConnectTarantool}
		_, err = taran.Ping()
		L.PanicIf(err, `taran.Ping`)
		mTest.Migrate(taran)
	}

	log.Println(`start test`)
	m.Run()
}

var runOnce = map[string]bool{}

func done() bool {
	caller := L.CallerInfo(2).FuncName
	//log.Println(caller)
	if runOnce[caller] {
		return true
	}
	runOnce[caller] = true
	return false
}

func timing() func(...int64) {
	start := time.Now()
	return func(v ...int64) {
		dur := time.Since(start)
		divisor := int64(total)
		if len(v) == 1 {
			divisor *= v[0]
		}
		// BenchmarkInsertS_Taran_ORM-32              10000             48616 ns/op            0.49 s
		fmt.Printf(`%-36s %11d %17d ns/op   %15.2f s`+"\n",
			fmt.Sprintf("%s-%d", L.CallerInfo(2).FuncName, cores),
			divisor,
			dur.Nanoseconds()/divisor,
			dur.Seconds(),
		)
	}
}
func idsToFetch(i uint64) []uint64 {
	return []uint64{(i + 100) % total, (i + 10) % total, (i + 1) % total, i % total}
}

func idsToFetchAny(i uint64) []any {
	return []any{(i + 100) % total, (i + 10) % total, (i + 1) % total, i % total}
}
