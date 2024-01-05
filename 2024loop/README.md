
# PostgreSQL, CockroachDB, Tarantool

## Test Env

- cockroachdb 23.1.13
- postgresql 16.1-1.pgdg120+1
- tarantool 2.11.2

```
goos: linux
goarch: amd64
```

## Usage

```shell
./clean-start.sh
go test -bench=Taran -benchmem .
go test -bench=Cockroach -benchmem .
go test -bench=Postgre -benchmem .
```

## Test Result

```
# Vinyl Engine
BenchmarkInsertS_Taran_ORM-32           100000   49350 ns/op     4.94 s
BenchmarkInsertS_Taran_ORM-32           100000   49351 ns/op     1144 B/op    29 allocs/op
BenchmarkUpdate_Taran_ORM-32            200000     148 ns/op     0.03 s
BenchmarkGetAllStruct_Taran_SQL-32        1417  842160 ns/op   936610 B/op  5734 allocs/op
BenchmarkGetAllStruct_Taran_ORM-32        3518  329884 ns/op   233959 B/op  4714 allocs/op
BenchmarkGetAllArray_Taran_ORM-32         3669  332334 ns/op   157568 B/op  4702 allocs/op
BenchmarkGetAllMap_Taran_SQL-32           1520  843275 ns/op  1248582 B/op  6732 allocs/op
BenchmarkGetOneStruct_Taran_SQL-32      105055   11961 ns/op     2569 B/op    57 allocs/op
BenchmarkGetOneMap_Taran_SQL-32         113568   10576 ns/op     2561 B/op    57 allocs/op
BenchmarkGetOneStruct_Taran_ORM-32      237246    5079 ns/op     1146 B/op    26 allocs/op
BenchmarkGetWhereInStruct_Taran_ORM-32   43108   23268 ns/op     3310 B/op    79 allocs/op
BenchmarkGetWhereInArray_Taran_ORM-32    50367   24336 ns/op     3254 B/op    79 allocs/op
BenchmarkGetLoopStruct_Taran_SQL-32      62161   18716 ns/op     4464 B/op   100 allocs/op

# MemTX Engine
BenchmarkInsertS_Taran_ORM-32           100000   48294 ns/op     4.83 s
BenchmarkInsertS_Taran_ORM-32           100000   48295 ns/op     1144 B/op    29 allocs/op
BenchmarkUpdate_Taran_ORM-32            200000     167 ns/op     0.03 s
BenchmarkGetAllStruct_Taran_SQL-32        2557  521430 ns/op   936605 B/op  5734 allocs/op
BenchmarkGetAllStruct_Taran_ORM-32       18490   66646 ns/op   233956 B/op  4714 allocs/op
BenchmarkGetAllArray_Taran_ORM-32        19908   61476 ns/op   157559 B/op  4702 allocs/op
BenchmarkGetAllMap_Taran_SQL-32           2262  534381 ns/op  1248570 B/op  6732 allocs/op
BenchmarkGetOneStruct_Taran_SQL-32      136198    8364 ns/op     2546 B/op    56 allocs/op
BenchmarkGetOneMap_Taran_SQL-32         150036    8400 ns/op     2539 B/op    56 allocs/op
BenchmarkGetOneStruct_Taran_ORM-32      250974    4436 ns/op     1146 B/op    26 allocs/op
BenchmarkGetWhereInStruct_Taran_SQL-32   71587   17438 ns/op     3311 B/op    79 allocs/op
BenchmarkGetWhereInArray_Taran_SQL-32    74480   16876 ns/op     3255 B/op    79 allocs/op
BenchmarkGetLoopStruct_Taran_ORM-32      69703   17131 ns/op     4472 B/op   100 allocs/op

# Postgres
benchmarkInsertPgx-32                   100000   49230 ns/op     4.92 s
BenchmarkInsert_Postgres_Pgx-32         100000   49231 ns/op      244 B/op     8 allocs/op
benchmarkUpdatePgx-32                   200000   48958 ns/op     9.79 s
BenchmarkUpdate_Postgres_Pgx-32         200000   48959 ns/op      200 B/op     7 allocs/op
BenchmarkGetAllStruct_Postgres_Pgx-32    23229   51164 ns/op    58490 B/op  2994 allocs/op
BenchmarkGetOneStruct_Postgres_Pgx-32   143138    8386 ns/op      578 B/op    13 allocs/op
BenchmarkGetWhereIn_Postgres_Pgx-32     146140    8432 ns/op      673 B/op    18 allocs/op
BenchmarkGetLoop_Postgres_Pgx-32         32282   37986 ns/op     4634 B/op   109 allocs/op

# Cockroach
benchmarkInsertPgx-32                   100000   97786 ns/op     9.78 s
BenchmarkInsert_Cockroach_Pgx-32        100000   97787 ns/op      242 B/op     8 allocs/op
benchmarkUpdatePgx-32                   200000  251920 ns/op     50.38 s
BenchmarkUpdate_Cockroach_Pgx-32        200000  251921 ns/op      202 B/op     7 allocs/op
BenchmarkGetAllStruct_Cockroach_Pgx-32   15722   79935 ns/op    58397 B/op  2947 allocs/op
BenchmarkGetOneStruct_Cockroach_Pgx-32   55669   20172 ns/op      579 B/op    13 allocs/op
BenchmarkGetWhereIn_Cockroach_Pgx-32     48294   25801 ns/op      674 B/op    18 allocs/op
BenchmarkGetLoop_Cockroach_Pgx-32        19368   59906 ns/op     5561 B/op   117 allocs/op

```


## Conclusion

Tarantool fastest for update, get single row use-case.
Postgres with pgx fastest for get multi-row use-case.
Cockroach with pgx slowest for update use-case.
`WHERE IN` always faster than loop `WHERE =`.
