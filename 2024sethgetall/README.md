
# Result

```
Redis 7

SET 10k user session  818.03 ms
12224 rps
GET 10k 20x user session 18023.34 ms
11097 rps, failed: 19
DEL 10k user session 704.50 ms

HSET 10k user session 1673.44 ms
5976 rps
HGETALL 10k 20x user session 18989.75 ms
10532 rps, failed: 23
DEL 10k user session 687.93 ms   

Tarantool 2.10

INSERT 10k user session 1209.11 ms
8271 rps
SELECT 10k 20x user session 19216.86 ms
10408 rps, failed: 0
```