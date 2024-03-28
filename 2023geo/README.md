
# Result (16-threads)

```
-- REDIS (by default indexed)

INSERTED 100K points: ok 100000 (100.0%) in 0.6 sec, 174372.6 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 33382 (16.7%) in 50.0 sec, 667.6 rps, points 16691000, 500.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 0.1 sec, 77609.4 rps, ERR: 0

RAM: 21 MB

-- REDIS (nearest first)

INSERTED 100K points: ok 100000 (100.0%) in 0.6 sec, 168513.6 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 30676 (15.3%) in 50.0 sec, 613.5 rps, points 15338000, 500.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 0.1 sec, 55088.4 rps, ERR: 0

-- KEYDB (16 threads)

INSERTED 100K points: ok 100000 (100.0%) in 1.3 sec, 79842.2 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 8382 (4.2%) in 50.0 sec, 167.6 rps, points 4191000, 500.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 0.1 sec, 50415.5 rps, ERR: 0

RAM: 36 MB  
  
-- POSTGRES (without index)

INSERTED 100K points: ok 100000 (100.0%) in 9.4 sec, 10639.7 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 200000 (100.0%) in 38.9 sec, 5136.6 rps, points 100000000, 500.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 0.4 sec, 11189.0 rps, ERR: 0

RAM: 26 MB

-- POSTGRES (with index)

INSERTED 100K points: ok 100000 (100.0%) in 9.2 sec, 10917.4 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 200000 (100.0%) in 37.5 sec, 5334.7 rps, points 100000000, 500.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 0.5 sec, 10978.2 rps, ERR: 0

-- POSTGRES (with index and distance ordering)

INSERTED 100K points: ok 100000 (100.0%) in 10.5 sec, 9523.7 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 7778 (3.9%) in 50.0 sec, 155.6 rps, points 3889000, 500.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 0.6 sec, 7954.3 rps, ERR: 0

-- TARANTOOL

INSERTED 100K points: ok 100000 (100.0%) in 0.8 sec, 129114.1 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 200000 (100.0%) in 29.8 sec, 6718.0 rps, points 100000000, 500.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 0.0 sec, 150719.2 rps, ERR: 0

-- TARATOOL (distance calculated on backend)

INSERTED 100K points: ok 100000 (100.0%) in 0.8 sec, 126801.9 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 200000 (100.0%) in 35.0 sec, 5716.3 rps, points 100000000, 500.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 0.0 sec, 137424.8 rps, ERR: 0

-- TYPESENSE (default 10 limit)

INSERTED 100K points: ok 98960 (100.0%) in 94.7 sec, 1044.5 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 34150 (17.1%) in 50.0 sec, 683.0 rps, points 341500, 10.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 9.6 sec, 523.5 rps, ERR: 0

-- TYPESENSE (250 limit, cannot be 500)

INSERTED 100K points: ok 98960 (100.0%) in 96.7 sec, 1023.8 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 2177 (1.1%) in 50.0 sec, 43.5 rps, points 544250, 250.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 9.7 sec, 515.2 rps, ERR: 0

-- MEILISEARCH (always zero result)

INSERTED 100K points: ok 98960 (100.0%) in 365.0 sec, 271.2 rps, ERR: 0
SEARCHED_RADIUS 200K points: ok 200000 (100.0%) in 4.8 sec, 41675.9 rps, points 0, 0.0 points/req ERR: 0
MOVING 5K points: ok 5000 (100.0%) in 8.8 sec, 569.7 rps, ERR: 0


```