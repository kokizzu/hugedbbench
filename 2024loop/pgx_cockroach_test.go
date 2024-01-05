package loopVsWhereIn

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgxCockroach *pgxpool.Pool

//
//func BenchmarkInsert_Cockroach_Pgx(b *testing.B) {
//	benchmarkInsertPgx(b, pgxCockroach)
//}
//
//func BenchmarkUpdate_Cockroach_Pgx(b *testing.B) {
//	benchmarkUpdatePgx(b, pgxCockroach)
//}
//
//func BenchmarkGetAllStruct_Cockroach_Pgx(b *testing.B) {
//	benchmarkGetAllStructPgx(b, pgxCockroach)
//}
//
//func BenchmarkGetRowS_Cockroach_Pgx(b *testing.B) {
//	benchmarkGetOneStructPgx(b, pgxCockroach)
//}
//
//func BenchmarkGetWhereIn_Cockroach_Pgx(b *testing.B) {
//	benchmarkWhereInPgx(b, pgxCockroach)
//}
