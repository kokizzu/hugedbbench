
# Benchmarking CockroachDB on different type of disk

| Disk Type       | InsDur | UpdDur | SelDur | ManyDur | InsQps | UpdQps | SelQps | ManyRow/s | ManyQps |
|-----------------|--------|--------|--------|---------|--------|--------|--------|-----------|---------|
| TMPFS (RAM)     |    1.3 |    2.1 |    4.9 |     1.5 |  31419 |  19275 | 81274  |   8194872 |   20487 |
| NVME DA 1TB     |    2.7 |    3.7 |    5.0 |     1.5 |  15072 |  10698 | 80558  |   8019435 |   20048 |
| NVMe Team 1TB   |    3.8 |    3.7 |    4.9 |     1.5 |  10569 |  10678 | 81820  |   8209889 |   20524 |
| SSD GALAX 250GB |    8.0 |    7.1 |    5.0 |     1.5 |   4980 |   5655 | 79877  |   7926162 |   19815 |
| HDD WD 8TB      |   32.1 |   31.7 |    4.9 |     3.9 |   1244 |   1262 | 81561  |   3075780 |    7689 |

- Dur in seconds
- Qps = queries per seconds, for Insert and Update equal to records per second since the insert/update done per 1 row
- Row/s number of rows scanned per second


## 0. TMPFS (RAM)

```
sudo mount -t tmpfs -o size=2G tmpfs tmp1

CockroachDB InsertOne: 1.3 sec | 31419.8 rec/s
CockroachDB Count 11.847716ms 40000
CockroachDB UpdateOne: 2.1 sec | 19275.9 rec/s
CockroachDB Count 32.326818ms 40000
CockroachDB SelectOne: 4.9 sec | 81274 rec/s
CockroachDB Select400: 1.5 sec | 8194872.5 rec/s | 20487.2 queries/s

sudo umount tmp1
```

## 1. NVMe DigitalAlliance 1 TB

```
CockroachDB InsertOne: 2.7 sec | 15072.3 rec/s
CockroachDB Count 14.90763ms 40000
CockroachDB UpdateOne: 3.7 sec | 10698.0 rec/s
CockroachDB Count 26.932609ms 40000
CockroachDB SelectOne: 5.0 sec | 80558 rec/s
CockroachDB Select400: 1.5 sec | 8019435.9 rec/s | 20048.6 queries/s
```

## 2. NVMe Team 512 GB

```
CockroachDB InsertOne: 3.8 sec | 10569.6 rec/s
CockroachDB Count 15.369907ms 40000
CockroachDB UpdateOne: 3.7 sec | 10678.3 rec/s
CockroachDB Count 23.200234ms 40000
CockroachDB SelectOne: 4.9 sec | 81820 rec/s
CockroachDB Select400: 1.5 sec | 8209889.1 rec/s | 20524.7 queries/s
```

## 3. SSD Galax 250 GB

```
CockroachDB InsertOne: 8.0 sec | 4980.8 rec/s
CockroachDB Count 19.208392ms 40000
CockroachDB UpdateOne: 7.1 sec | 5655.3 rec/s
CockroachDB Count 40.879448ms 40000
CockroachDB SelectOne: 5.0 sec | 79877 rec/s
CockroachDB Select400: 1.5 sec | 7926162.5 rec/s | 19815.4 queries/s
```

## 4. HDD WD 8 TB

```
CockroachDB InsertOne: 32.1 sec | 1244.2 rec/s
CockroachDB Count 15.54701ms 40000
CockroachDB UpdateOne: 31.7 sec | 1262.0 rec/s
CockroachDB Count 32.782333ms 40000
CockroachDB SelectOne: 4.9 sec | 81561 rec/s
CockroachDB Select400: 3.9 sec | 3075780.1 rec/s | 7689.5 queries/s
```

CockroachDB version: v22.1.8 @ 2022/09/29 14:21:51 (go1.17.11).
Disk Usage by Database: 1.6G
