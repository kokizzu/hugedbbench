package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/twmb/franz-go/pkg/kgo"
	"sync"
	"sync/atomic"
	"time"
)

const PRODUCERS = 1000
const MSGS = 2000 // x PRODUCERS
const CONSUMERS = 100
const TOPIC = `foo`
const PROGRESS = 10000

/*
pre-config:

/etc/sysctl.conf:
fs.aio-max-nr = 4194304

sudo sysctl -p
*/

// docker-compose -f docker-compose-single.yaml up --remove-orphans
// docker exec -it $(docker ps | grep redpanda1 | cut -d ' ' -f 1) rpk topic create foo

// docker-compose -f docker-compose-multi.yaml up --remove-orphans
// docker exec -it $(docker ps | grep redpanda1 | cut -d ' ' -f 1) rpk topic create foo

func main() {
	startBenchmark := time.Now()
	//seeds := []string{"localhost:9092"}
	seeds := []string{"localhost:9092", "localhost:9093", "localhost:9094"}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("group1"),
		kgo.ConsumeTopics(TOPIC),
		kgo.AllowAutoTopicCreation(),
	)
	L.PanicIf(err, `kgo.NewClient`)
	defer cl.Close()

	wgConsume := &sync.WaitGroup{}
	wgConsume.Add(PRODUCERS * MSGS) // includes consuming
	failConsume := int64(0)
	doubleConsume := int64(0)
	durConsume := int64(0)
	maxLatency := int64(0)
	consumed := int64(0)
	consume := sync.Map{}

	go func() {
		for z := 0; z < CONSUMERS; z++ {
			go func(z int) {
				//fmt.Println(`Consumer spawned`, z)
				for {
					fetches := cl.PollFetches(context.Background())
					if errs := fetches.Errors(); len(errs) > 0 {
						atomic.AddInt64(&failConsume, int64(len(errs)))
						L.Print(errs)
					}
					iter := fetches.RecordIter()
					for !iter.Done() {
						record := iter.Next()
						m := M.SX{}
						err := json.Unmarshal(record.Value, &m)
						if err != nil {
							atomic.AddInt64(&failConsume, 1)
							L.Print(err)
						}
						if _, exists := consume.LoadOrStore(m.GetInt(`idx`), struct{}{}); !exists {
							dur := time.Now().UnixNano() - m.GetInt(`when`)
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
			}(z)
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
				record := &kgo.Record{Topic: TOPIC, Value: []byte(M.SX{
					`when`: time.Now().UnixNano(),
					`idx`:  z*MSGS + y,
				}.ToJson())}
				cl.Produce(context.Background(), record, func(_ *kgo.Record, err error) {
					defer wgProduce.Done()
					//fmt.Println(record)
					if err != nil {
						atomic.AddInt64(&failProduce, 1)
						L.Print(err)
						return
					}
					if atomic.AddInt64(&produced, 1)%PROGRESS == 0 {
						//fmt.Print("P")
					}
				})
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
	fmt.Println(`MaxLatency (ms): `, maxLatency/1000/1000)
	fmt.Println(`AvgLatency (ms): `, durConsume/PRODUCERS/MSGS/1000/1000)
	fmt.Println(`Total (s)`, time.Since(startBenchmark))
}
