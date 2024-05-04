
# Result Synchronous

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

MemTx Engine
  
INSERT 10k user session 1209.11 ms
8271 rps
SELECT 10k 20x user session 19216.86 ms
10408 rps, failed: 0

Vinyl Engine

INSERT 10k user session 1140.26 ms
8770 rps
SELECT 10k 20x user session 18461.86 ms
10833 rps, failed: 0
```

# Result Asynchronous/100 Thread

```

Redis 7  

SET 10k user session, 100 thread 42.63 ms
234579 rps
GET 10k 20x user session 718.85 ms
278221 rps, failed: 0
DEL 10k user session 753.18 ms

HSET+TTL 10k user session, 100 thread 68.75 ms
145457 rps
HGETALL 10k 20x user session 748.57 ms
267177 rps, failed: 0
DEL 10k user session 725.20 ms
  

Tarantool 2.11
  
Vinyl Engine  
 
INSERT 10k user session, 100 thread 82.82 ms
120737 rps
SELECT 10k 20x user session 1361.17 ms
146933 rps, failed: 0


Tarantool 3.1

Vinyl Engine

INSERT 10k user session, 100 thread 39.13 ms
255569 rps

SELECT 10k 20x user session 567.16 ms
352632 rps, failed: 0
        
```