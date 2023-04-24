
# Result (16-threads)

```
-- REDIS

INSERTED 100K points: ok 100000 (100.0%) in 0.6 sec, 174372.6 rps, ERR: {0}
SEARCHED_RADIUS 200K points: ok 33382 (16.7%) in 50.0 sec, 667.6 rps, points 16691000, 500.0 points/req ERR: {0}
MOVING 5K points: ok 5000 (100.0%) in 0.1 sec, 77609.4 rps, ERR: {0}

RAM: 21 MB

-- POSTGRES

INSERTED 100K points: ok 100000 (100.0%) in 9.4 sec, 10639.7 rps, ERR: {0}
SEARCHED_RADIUS 200K points: ok 200000 (100.0%) in 38.9 sec, 5136.6 rps, points 100000000, 500.0 points/req ERR: {0}
MOVING 5K points: ok 5000 (100.0%) in 0.4 sec, 11189.0 rps, ERR: {0}

RAM: 26 MB

```