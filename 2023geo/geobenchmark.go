package geo

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alitto/pond"
	"github.com/kpango/fastime"
	"go.uber.org/atomic"
)

const lat = 1.290270
const long = 103.851959
const DistanceDeg = 0.04   // est 4km (2km to 4 direction)
const DegToMeter = 111_000 // should be 111,139
const distanceBoxMeter = DegToMeter * DistanceDeg
const maxResult = 500
const workercount = 16
const workerCap = 1024

// Insert100kPoints insert 100K points
func Insert100kPoints(insertFunc func(lat float64, long float64, id uint64) error) {
	start := time.Now()

	rand.Seed(99) // make it deterministic
	errCounter := atomic.Uint64{}
	okCounter := atomic.Uint64{}
	const total = 100_000
	pool := pond.New(workercount, workerCap)
	for id := uint64(1); id <= total; id++ {
		id := id
		pool.Submit(func() {
			lat := (rand.Float64()-0.5)*0.2 + lat
			long := (rand.Float64()-0.5)*0.2 + long
			if insertFunc(lat, long, id) != nil {
				errCounter.Inc()
			} else {
				okCounter.Inc()
			}
		})
	}
	loopEnd := time.Now()
	pool.StopAndWaitFor(30*time.Second - loopEnd.Sub(start))
	duration := time.Since(start).Seconds()
	okCount := okCounter.Load()
	rps := float64(okCount) / duration
	percent := float64(okCount) * 100 / float64(okCount)
	fmt.Printf("INSERTED 100K points: ok %d (%.1f%%) in %.1f sec, %.1f rps, ERR: %d\n",
		okCount, percent, duration, rps, errCounter.Load())
}

// SearchRadius200k 200K times searching for points
func SearchRadius200k(searchFunc func(lat float64, long float64, boxMeter float64, maxResult int64) (uint64, error)) {
	start := time.Now()
	rand.Seed(88) // make it deterministic
	errCounter := atomic.Uint64{}
	okCounter := atomic.Uint64{}
	totalFound := atomic.Uint64{}
	const total = 200_000
	pool := pond.New(workercount, workerCap)
	for z := 0; z < total; z++ {
		pool.Submit(func() {
			lat := (rand.Float64()-0.5)*0.2 + lat
			long := (rand.Float64()-0.5)*0.2 + long
			if pts, err := searchFunc(lat, long, distanceBoxMeter, maxResult); err != nil {
				errCounter.Inc()
			} else {
				totalFound.Add(pts)
				okCounter.Inc()
			}
		})
		if z%100 == 1 { // every 100 request, check if time nearly exceeded
			if fastime.Now().Sub(start).Seconds() >= 49 {
				break
			}
		}
	}
	loopEnd := time.Now()
	pool.StopAndWaitFor(50*time.Second - loopEnd.Sub(start))
	duration := time.Since(start).Seconds()
	okCount := okCounter.Load()
	rps := float64(okCount) / duration
	percent := float64(okCount) * 100 / float64(total)
	found := totalFound.Load()
	ptsPerReq := float64(found) / float64(okCount)
	fmt.Printf("SEARCHED_RADIUS 200K points: ok %d (%.1f%%) in %.1f sec, %.1f rps, points %d, %.1f points/req ERR: %d\n",
		okCount, percent, duration, rps, found, ptsPerReq, errCounter.Load())
}

// MovingPoint moving first 100 point to another location, each 50 times
func MovingPoint(movingFunc func(lat float64, long float64, id uint64) error) {
	start := time.Now()
	rand.Seed(77) // make it deterministic
	errCounter := atomic.Uint64{}
	okCounter := atomic.Uint64{}
	const total = 5000
	pool := pond.New(workercount, workerCap)
	for id := uint64(1); id <= total; id++ {
		id := id
		pool.Submit(func() {
			lat := (rand.Float64()-0.5)*0.2 + lat
			long := (rand.Float64()-0.5)*0.2 + long
			if movingFunc(lat, long, id%100) != nil {
				errCounter.Inc()
			} else {
				okCounter.Inc()
			}
		})
	}
	loopEnd := time.Now()
	pool.StopAndWaitFor(10*time.Second - loopEnd.Sub(start))
	duration := time.Since(start).Seconds()
	okCount := okCounter.Load()
	rps := float64(okCount) / duration
	percent := float64(okCount) * 100 / float64(okCount)
	fmt.Printf("MOVING 5K points: ok %d (%.1f%%) in %.1f sec, %.1f rps, ERR: %d\n",
		okCount, percent, duration, rps, errCounter.Load())
}
