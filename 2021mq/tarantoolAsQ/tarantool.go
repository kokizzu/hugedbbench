package main

import (
	"fmt"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/tarantool/go-tarantool"
	"hugedbbench/2021mq/tarantoolAsQ/mFoo"
	"hugedbbench/2021mq/tarantoolAsQ/mFoo/rqFoo"
	"hugedbbench/2021mq/tarantoolAsQ/mFoo/wcFoo"
	"sync"
	"sync/atomic"
	"time"
)

func ConnectTarantool() *tarantool.Connection {
	hostPort := fmt.Sprintf(`%s:%s`,
		`127.0.0.1`,
		`3301`,
	)
	taran, err := tarantool.Connect(hostPort, tarantool.Opts{
		User: `guest`,
		Pass: ``,
	})
	L.PanicIf(err, `tarantool.Connect `+hostPort)
	return taran
}

// then use it like this:

const PRODUCERS = 1000
const MSGS = 2000 // x PRODUCERS
const CONSUMERS = 100
const TOPIC = `foo`
const PROGRESS = 10000

/*
using sequence as counter  https://www.tarantool.io/en/doc/latest/book/box/data_model/#index-box-sequence

TARANTOOL version: 2.8.2
standard database, with 1 sec pooling delay
Memtx engine, because Vinyl engine not true write-serializable
probably must use auto_increment

=== tarantool single (using manual counter, not sequence, unstable):

FailProduce:  0
FailConsume:  0
DoubleConsume:  0
Produced (ms):  7380
MaxLatency (ms):  79424
AvgLatency (ms):  77981
Total (s) 13.087616048s

FailProduce:  0
FailConsume:  0
DoubleConsume:  0
Produced (ms):  6064
MaxLatency (ms):  5533
AvgLatency (ms):  4580
Total (s) 10.213670874s
*/

// docker-compose -f docker-compose.yaml up --remove-orphans

// cleanup
// docker exec -it $(docker ps | grep tarantool | cut -f 1 -d ' ') tarantoolctl connect 3301
//   box.space.foo:truncate()
//   box.space.foo:drop() # or drop

// check/verify
// docker exec -it $(docker ps | grep taranto | cut -d ' ' -f 1) tarantoolctl connect 3301
//   box.space.foo:count() # should be 2000000

func main() {
	startBenchmark := time.Now()
	tt := &Tt.Adapter{Connection: ConnectTarantool(), Reconnect: ConnectTarantool}
	defer tt.Close()

	if !tt.CreateSpace(TOPIC, Tt.Vinyl) {
		panic(`tt.CreateSpace`)
	}
	if !tt.UpsertTable(TOPIC, mFoo.TarantoolTables[mFoo.TableFoo]) {
		panic(`tt.UpsertTable`)
	}
	if !tt.TruncateTable(TOPIC) {
		panic(`tt.TruncateTable`)
	}
	
	wgConsume := &sync.WaitGroup{}
	wgConsume.Add(PRODUCERS * MSGS) // includes consuming
	failConsume := int64(0)
	doubleConsume := int64(0)
	durConsume := int64(0)
	maxLatency := int64(0)
	consumed := int64(0)
	consume := sync.Map{}
	
	go func() {
		lastFetch := int64(0)
		rq := rqFoo.NewFoo(tt)
		for {
			rows := rq.FindGreaterThan(lastFetch, 10000)
			if len(rows) == 0 {
				time.Sleep(time.Second)
			}
			for _, row := range rows {
				idx := int64(row.Id)
				createdAt := int64(row.When)
				if lastFetch < idx {
					lastFetch = idx
				}
				// single processing
				if _, exists := consume.LoadOrStore(idx, struct{}{}); !exists {
					dur := time.Now().UnixNano()/1000/1000 - createdAt/1000/1000 // millisecond
					atomic.AddInt64(&durConsume, dur) // bottleneck, TODO: change to channel
					for {
						old := maxLatency
						if dur <= old {
							break
						}
						if atomic.CompareAndSwapInt64(&maxLatency, old, dur) {
							break
						}
					}
					if atomic.AddInt64(&consumed, 1)%PROGRESS == 0 {
						//fmt.Print("C")
					}
					wgConsume.Done()
				} else {
					atomic.AddInt64(&doubleConsume, 1)
				}
			}
		}
	}()
	
	wgProduce := &sync.WaitGroup{}
	wgProduce.Add(PRODUCERS * MSGS)
	failProduce := int64(0)
	durProduce := int64(0)
	produced := int64(0)
	
	startProduce := time.Now().UnixNano()
	for z := 0; z < PRODUCERS; z++ {
		go func(z int) {
			//fmt.Println(`Producer spawned`, z)
			for y := 0; y < MSGS; y++ {
				wc := wcFoo.NewFooMutator(tt)
				wc.When = uint64(time.Now().UnixNano())
				if !wc.DoInsert() {
					atomic.AddInt64(&failProduce, 1)
					return
				}
				if atomic.AddInt64(&produced, 1)%PROGRESS == 0 {
					//fmt.Print("P")
				}
				wgProduce.Done()
			}
		}(z)
	}
	
	wgProduce.Wait()
	durProduce = time.Now().UnixNano() - startProduce
	wgConsume.Wait()
	
	fmt.Println(`FailProduce: `, failProduce)
	fmt.Println(`FailConsume: `, failConsume)
	fmt.Println(`DoubleConsume: `, doubleConsume)
	fmt.Println(`Produced (ms): `, durProduce/1000/1000)
	fmt.Println(`MaxLatency (ms): `, maxLatency)
	fmt.Println(`AvgLatency (ms): `, durConsume/PRODUCERS/MSGS)
	fmt.Println(`Total (s)`, time.Since(startBenchmark))
}
