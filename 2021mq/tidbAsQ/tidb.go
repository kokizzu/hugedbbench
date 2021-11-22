package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kokizzu/gotro/L"
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
 not ok, since TiDB inserts not true serializable like in old SQL databases
*/

// docker-compose -f docker-compose.yaml up --remove-orphans

func main() {
	startBenchmark := time.Now()
	//seeds := []string{"localhost:9092"}

	bg := context.Background() // shared context
	myUrl := "%s:%s@tcp(%s:%d)/%s"
	myUrl = fmt.Sprintf(myUrl,
		`root`,
		``, // empty password
		`127.0.0.1`,
		4000,
		`test`,
	)

	conn, err := sql.Open("mysql", myUrl)
	L.PanicIf(err, `cannot connect db`)
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(1024)
	conn.SetMaxIdleConns(1024)
	defer conn.Close()

	_, err = conn.ExecContext(bg, `CREATE TABLE IF NOT EXISTS foo(idx BIGINT PRIMARY KEY AUTO_INCREMENT, createdAt BIGINT)`)
	L.PanicIf(err, `failed create table bar1`)

	_, err = conn.ExecContext(bg, `TRUNCATE TABLE foo`)
	L.PanicIf(err, `failed truncate table bar1`)

	wgConsume := &sync.WaitGroup{}
	wgConsume.Add(PRODUCERS * MSGS) // includes consuming
	failConsume := int64(0)
	doubleConsume := int64(0)
	durConsume := int64(0)
	maxLatency := int64(0)
	consumed := int64(0)
	consume := sync.Map{}

	go func() {
		lastFetch := int64(-1)
		for {
			var idx, createdAt int64
			rows, err := conn.QueryContext(bg, `SELECT idx, createdAt FROM foo WHERE idx>? ORDER BY idx LIMIT 1000`, lastFetch)
			if err != nil {
				if err == sql.ErrNoRows {
					time.Sleep(time.Second)
				} else {
					L.Print(err)
				}
			}
			for rows.Next() {
				err := rows.Scan(&idx, &createdAt)
				if err != nil {
					atomic.AddInt64(&failConsume, 1)
					L.Print(err)
				}
				if lastFetch < idx {
					lastFetch = idx
				}
				// single processing
				if _, exists := consume.LoadOrStore(idx, struct{}{}); !exists {
					dur := time.Now().UnixNano() - createdAt
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
						fmt.Print("C")
					}
					wgConsume.Done()
				} else {
					atomic.AddInt64(&doubleConsume, 1)
				}
			}
			rows.Close()
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
				_, err = conn.ExecContext(bg, `INSERT INTO foo(createdAt) VALUES(?)`, time.Now().UnixNano())
				if err != nil {
					atomic.AddInt64(&failProduce, 1)
					L.Print(err)
					return
				}
				if atomic.AddInt64(&produced, 1)%PROGRESS == 0 {
					fmt.Print("P")
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
	fmt.Println(`MaxLatency (ms): `, maxLatency/1000/1000)
	fmt.Println(`AvgLatency (ms): `, durConsume/PRODUCERS/MSGS/1000/1000)
	fmt.Println(`Total (s)`, time.Since(startBenchmark))
}
