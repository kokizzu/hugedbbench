package main

import (
	"encoding/json"
	"fmt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const PRODUCERS = 100 //* 0.5 // some publishers timed out
const MSGS = 20000    // x PRODUCERS
const CONSUMERS = 100 //* 0.1 // dunno why, 56th-86th consumer always timeout
const TOPIC = `foo`
const PROGRESS = 10000
const WILDCARD = `.*`

// docker-compose -f docker-compose-single.yaml up --remove-orphans
// docker-compose -f docker-compose-multi.yaml up --remove-orphans

// https://shijuvar.medium.com/building-distributed-event-streaming-systems-in-go-with-nats-jetstream-3938e6dc7a13
func main() {
	startBenchmark := time.Now()
	nc, err := nats.Connect(nats.DefaultURL)
	L.PanicIf(err, `nats.Connect`)
	js, err := nc.JetStream()
	L.PanicIf(err, `nc.JetStream`)
	defer nc.Close()

	// create stream
	stream, err := js.StreamInfo(TOPIC)
	L.IsError(err, `js.StreamInfo`)
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", TOPIC, TOPIC)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     TOPIC,
			Subjects: []string{TOPIC},
		})
		L.IsError(err, `js.AddStream`)
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
		for z := 0; z < CONSUMERS; z++ {
			go func(z int) {
				//fmt.Println(`Consumer spawned`, z)
				_, err := js.Subscribe(TOPIC, func(msg *nats.Msg) {
					//atomic.AddInt64(&failConsume, int64(len(errs)))
					//L.Print(errs)
					m := M.SX{}
					err := json.Unmarshal(msg.Data, &m)
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
						_ = msg.Ack()
						wgConsume.Done()
					} else {
						atomic.AddInt64(&doubleConsume, 1)
					}
				})
				L.PanicIf(err, `js.Subscribe %d`, z)
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
				_, err := js.Publish(TOPIC, []byte(M.SX{
					`when`: time.Now().UnixNano(),
					`idx`:  z*MSGS + y,
				}.ToJson()))
				if err != nil {
					atomic.AddInt64(&failProduce, 1)
					L.Print(err, z)
					return
				}
				wgProduce.Done()
				if atomic.AddInt64(&produced, 1)%PROGRESS == 0 {
					//fmt.Print("P")
				}
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
