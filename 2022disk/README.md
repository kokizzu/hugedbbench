
# Benchmarking CockroachDB on different type of disk

## 0. TMPFS (RAM)

```
sudo mount -t tmpfs -o size=2G tmpfs tmp1

CockroachDB InsertOne: 1.3 sec | 31419.8 rec/s
CockroachDB Count 11.847716ms 40000
CockroachDB UpdateOne: 2.1 sec | 19275.9 rec/s
CockroachDB Count 32.326818ms 40000
CockroachDB SelectOne: 4.9 sec | 8127.4 rec/s
CockroachDB Select400: 1.5 sec | 8194872.5 rec/s | 20487.2 queries/s

sudo umount tmp1
```

## 1. NVMe DigitalAlliance 1 TB

```
CockroachDB InsertOne: 2.7 sec | 15072.3 rec/s
CockroachDB Count 14.90763ms 40000
CockroachDB UpdateOne: 3.7 sec | 10698.0 rec/s
CockroachDB Count 26.932609ms 40000
CockroachDB SelectOne: 5.0 sec | 8055.8 rec/s
CockroachDB Select400: 1.5 sec | 8019435.9 rec/s | 20048.6 queries/s
```

# 2. NVMe Team 512 GB

```
CockroachDB InsertOne: 3.8 sec | 10569.6 rec/s
CockroachDB Count 15.369907ms 40000
CockroachDB UpdateOne: 3.7 sec | 10678.3 rec/s
CockroachDB Count 23.200234ms 40000
CockroachDB SelectOne: 4.9 sec | 8182.0 rec/s
CockroachDB Select400: 1.5 sec | 8209889.1 rec/s | 20524.7 queries/s
```

# 3. SSD Galax 250 GB

```
CockroachDB InsertOne: 8.0 sec | 4980.8 rec/s
CockroachDB Count 19.208392ms 40000
CockroachDB UpdateOne: 7.1 sec | 5655.3 rec/s
CockroachDB Count 40.879448ms 40000
CockroachDB SelectOne: 5.0 sec | 7987.7 rec/s
CockroachDB Select400: 1.5 sec | 7926162.5 rec/s | 19815.4 queries/s
```

# 4. HDD WD 8 TB

```
CockroachDB InsertOne: 32.1 sec | 1244.2 rec/s
CockroachDB Count 15.54701ms 40000
CockroachDB UpdateOne: 31.7 sec | 1262.0 rec/s
CockroachDB Count 32.782333ms 40000
CockroachDB SelectOne: 4.9 sec | 8156.1 rec/s
CockroachDB Select400: 3.9 sec | 3075780.1 rec/s | 7689.5 queries/s
```