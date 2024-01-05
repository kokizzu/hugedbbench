package loopVsWhereIn

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgxPostgres *pgxpool.Pool

//func BenchmarkInsert_Postgres_Pgx(b *testing.B) {
//	benchmarkInsertPgx(b, pgxPostgres)
//}
//
//func BenchmarkUpdate_Postgres_Pgx(b *testing.B) {
//	benchmarkUpdatePgx(b, pgxPostgres)
//}
//
//func BenchmarkGetAllStruct_Postgres_Pgx(b *testing.B) {
//	benchmarkGetAllStructPgx(b, pgxPostgres)
//}
//
//func BenchmarkGetOneStruct_Postgres_Pgx(b *testing.B) {
//	benchmarkGetOneStructPgx(b, pgxPostgres)
//}
//
//func BenchmarkGetWhereIn_Postgres_Pgx(b *testing.B) {
//	benchmarkWhereInPgx(b, pgxPostgres)
//}
//
//func BenchmarkGetLoop_Postgres_Pgx(b *testing.B) {
//	benchmarkLoopPgx(b, pgxPostgres)
//}
