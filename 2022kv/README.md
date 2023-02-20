
# KV use-case benchmark

key-value use-case benchmark, including counters

benchmark: 
- 5s put / 1m recs x 100 goroutine
- 8s get / 4m recs x 1000 goroutine
- 2s inc/dec / 1m recs x 100 goroutine
- 2s del / 1m recs x 100 goroutine

products:
- tarantool memtx (in-memory, WAL)
- tarantool vinyl (disk-based)
- aerospike CE (in-memory)
- redis (in-memory, RDB+AOF)
- mongodb (disk-based)
- scylladb (disk-based)
- keydb (in-memory, RDB+AOF)
- dragonflydb (in-memory)
- tidb (disk-based)
- cockroachdb (disk-based)
- postgresql (disk-based)
- mysql (disk-based)
- singlestore (disk-based)
- yugabytedb (disk-based)
- icefiredb (in-memory)
- singlestore (disk-based)
- cratedb (disk-based)
- tendis (disk-based)

measurement:
- memory usage
- disk usage (if any)
- rps
- how easy to add new replica