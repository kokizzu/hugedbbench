
# Result (16-threads)

```
-- REDIS (by default indexed)

INSERTED 100K points: ok 100000 (100.0%) in 0.6 sec, 174372.6 rps, ERR: {0}
SEARCHED_RADIUS 200K points: ok 33382 (16.7%) in 50.0 sec, 667.6 rps, points 16691000, 500.0 points/req ERR: {0}
MOVING 5K points: ok 5000 (100.0%) in 0.1 sec, 77609.4 rps, ERR: {0}

RAM: 21 MB

-- REDIS (nearestt first)

INSERTED 100K points: ok 100000 (100.0%) in 0.6 sec, 168513.6 rps, ERR: {0}
SEARCHED_RADIUS 200K points: ok 30676 (15.3%) in 50.0 sec, 613.5 rps, points 15338000, 500.0 points/req ERR: {0}
MOVING 5K points: ok 5000 (100.0%) in 0.1 sec, 55088.4 rps, ERR: {0}

-- POSTGRES (without index)

INSERTED 100K points: ok 100000 (100.0%) in 9.4 sec, 10639.7 rps, ERR: {0}
SEARCHED_RADIUS 200K points: ok 200000 (100.0%) in 38.9 sec, 5136.6 rps, points 100000000, 500.0 points/req ERR: {0}
MOVING 5K points: ok 5000 (100.0%) in 0.4 sec, 11189.0 rps, ERR: {0}

RAM: 26 MB

-- POSTGRES (with index)

INSERTED 100K points: ok 100000 (100.0%) in 9.2 sec, 10917.4 rps, ERR: {0}
SEARCHED_RADIUS 200K points: ok 200000 (100.0%) in 37.5 sec, 5334.7 rps, points 100000000, 500.0 points/req ERR: {0}
MOVING 5K points: ok 5000 (100.0%) in 0.5 sec, 10978.2 rps, ERR: {0}

-- POSTGRES (with index and distance ordering)

INSERTED 100K points: ok 100000 (100.0%) in 10.5 sec, 9523.7 rps, ERR: {0}
SEARCHED_RADIUS 200K points: ok 7778 (3.9%) in 50.0 sec, 155.6 rps, points 3889000, 500.0 points/req ERR: {0}
MOVING 5K points: ok 5000 (100.0%) in 0.6 sec, 7954.3 rps, ERR: {0}


```