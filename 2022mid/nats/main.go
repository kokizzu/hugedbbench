package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/nats-io/nats.go"
)

func main() {
	tStr := `_` + I.ToS(time.Now().UnixNano())
	if len(os.Args) > 1 {
		app := fiber.New()

		mode := os.Args[1]
		switch mode {
		case `apiserver`:
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString(I.ToS(rand.Int63()) + tStr)
			})

		case `apiproxy`:
			// connect as request on request-reply

			const N = 8
			counter := uint32(0)
			ncs := [N]*nats.Conn{}
			mutex := sync.Mutex{}
			conn := func() *nats.Conn {
				idx := atomic.AddUint32(&counter, 1) % N
				nc := ncs[idx]
				if nc != nil {
					return nc
				}
				mutex.Lock()
				if ncs[idx] != nil {
					return ncs[idx]
				}
				nc, err := nats.Connect("127.0.0.1")
				L.PanicIf(err, `nats.Connect`)
				ncs[idx] = nc
				mutex.Unlock()
				return nc
			}

			defer func() {
				for _, nc := range ncs {
					if nc != nil {
						nc.Close()
					}
				}
			}()

			// handler
			app.Get("/", func(c *fiber.Ctx) error {
				msg, err := conn().Request("handle1", []byte(I.ToS(rand.Int63())), time.Second)
				if L.IsError(err, `nc.Request`) {
					return err
				}

				// Use the response
				return c.SendString(string(msg.Data))
			})
		default:
		}

		log.Println(mode + ` started ` + tStr)
		log.Fatal(app.Listen(":3000"))

	} else {
		// worker
		log.Println(`worker started ` + tStr)

		nc, err := nats.Connect("127.0.0.1")
		L.PanicIf(err, `nats.Connect`)
		defer nc.Close()

		// connect as reply on request-reply
		sub, err := nc.SubscribeSync("handle1")
		L.PanicIf(err, `nc.SubscribeSync`)

		// Wait for a message
		for {
			msg, err := sub.NextMsgWithContext(context.Background())
			L.PanicIf(err, `sub.NextMsgWithContext`)

			err = msg.Respond([]byte(string(msg.Data) + tStr))
			L.PanicIf(err, `msg.Respond`)
		}

		// TODO: try channel and callback version, also try reponse using goroutine
	}
}

/*
benchmark scenario:

###########################################################################

1. direct handling

go run main.go apiserver

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        4.3014 secs
  Slowest:      0.1061 secs
  Fastest:      0.0000 secs
  Average:      0.0011 secs
  Requests/sec: 232449.1716

  Total data:   38873797 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.011 [998810]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.021 [761]   |
  0.032 [18]    |
  0.042 [95]    |
  0.053 [2]     |
  0.064 [45]    |
  0.074 [34]    |
  0.085 [7]     |
  0.095 [0]     |
  0.106 [82]    |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0001 secs
  75% in 0.0019 secs
  90% in 0.0035 secs
  95% in 0.0041 secs
  99% in 0.0070 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.1061 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0623 secs
  resp wait:    0.0000 secs, 0.0000 secs, 0.0618 secs
  resp read:    0.0005 secs, 0.0000 secs, 0.0982 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

2. proxied with 1 worker, single nats

docker-compose up -d

go run main.go apiproxy

go run main.go

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        9.9526 secs
  Slowest:      0.0768 secs
  Fastest:      0.0002 secs
  Average:      0.0025 secs
  Requests/sec: 100461.5866

  Total data:   38872892 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.008 [999487]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.015 [62]    |
  0.023 [114]   |
  0.031 [47]    |
  0.038 [0]     |
  0.046 [47]    |
  0.054 [59]    |
  0.061 [1]     |
  0.069 [0]     |
  0.077 [37]    |


Latency distribution:
  10% in 0.0020 secs
  25% in 0.0022 secs
  50% in 0.0025 secs
  75% in 0.0027 secs
  90% in 0.0030 secs
  95% in 0.0033 secs
  99% in 0.0042 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0002 secs, 0.0768 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0525 secs
  resp wait:    0.0024 secs, 0.0001 secs, 0.0718 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0445 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

3. proxied with 2 worker, single nats

go run main.go # start another instance

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        13.1335 secs
  Slowest:      0.0894 secs
  Fastest:      0.0002 secs
  Average:      0.0033 secs
  Requests/sec: 76130.4079

  Total data:   38874172 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.009 [999387]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.018 [67]    |
  0.027 [68]    |
  0.036 [272]   |
  0.045 [0]     |
  0.054 [0]     |
  0.063 [55]    |
  0.072 [4]     |
  0.080 [0]     |
  0.089 [1]     |


Latency distribution:
  10% in 0.0027 secs
  25% in 0.0030 secs
  50% in 0.0033 secs
  75% in 0.0036 secs
  90% in 0.0040 secs
  95% in 0.0043 secs
  99% in 0.0052 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0002 secs, 0.0894 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0605 secs
  resp wait:    0.0032 secs, 0.0002 secs, 0.0321 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0334 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

4. proxied with 4 workers, single nats

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        19.9410 secs
  Slowest:      0.1065 secs
  Fastest:      0.0003 secs
  Average:      0.0051 secs
  Requests/sec: 50140.6288

  Total data:   38874052 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.011 [999551]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.022 [117]   |
  0.032 [69]    |
  0.043 [34]    |
  0.053 [0]     |
  0.064 [44]    |
  0.075 [29]    |
  0.085 [0]     |
  0.096 [0]     |
  0.106 [10]    |


Latency distribution:
  10% in 0.0038 secs
  25% in 0.0043 secs
  50% in 0.0050 secs
  75% in 0.0057 secs
  90% in 0.0063 secs
  95% in 0.0068 secs
  99% in 0.0082 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0003 secs, 0.1065 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0616 secs
  resp wait:    0.0050 secs, 0.0002 secs, 0.0401 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.1018 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

what if we limit CPU?

5. apiserver 2 core

GOMAXPROCS=2 go run main.go apiserver

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        5.4271 secs
  Slowest:      0.0879 secs
  Fastest:      0.0000 secs
  Average:      0.0014 secs
  Requests/sec: 184234.0106

  Total data:   38874151 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.009 [999568]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.018 [10]    |
  0.026 [123]   |
  0.035 [21]    |
  0.044 [22]    |
  0.053 [0]     |
  0.062 [74]    |
  0.070 [3]     |
  0.079 [0]     |
  0.088 [33]    |


Latency distribution:
  10% in 0.0005 secs
  25% in 0.0008 secs
  50% in 0.0013 secs
  75% in 0.0018 secs
  90% in 0.0023 secs
  95% in 0.0026 secs
  99% in 0.0037 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0879 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0597 secs
  resp wait:    0.0013 secs, 0.0000 secs, 0.0605 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0847 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

6. 2 core apiproxy, 1 worker 2 core, single nats

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        9.7066 secs
  Slowest:      0.0942 secs
  Fastest:      0.0001 secs
  Average:      0.0025 secs
  Requests/sec: 103007.4516

  Total data:   38873771 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.010 [999537]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.019 [12]    |
  0.028 [35]    |
  0.038 [123]   |
  0.047 [0]     |
  0.057 [45]    |
  0.066 [53]    |
  0.075 [0]     |
  0.085 [0]     |
  0.094 [49]    |


Latency distribution:
  10% in 0.0015 secs
  25% in 0.0019 secs
  50% in 0.0024 secs
  75% in 0.0029 secs
  90% in 0.0035 secs
  95% in 0.0039 secs
  99% in 0.0049 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.0942 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0606 secs
  resp wait:    0.0024 secs, 0.0001 secs, 0.0604 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0564 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

7. 2 core apiproxy, 2 worker 2 core, single nats

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        11.4240 secs
  Slowest:      0.0945 secs
  Fastest:      0.0001 secs
  Average:      0.0029 secs
  Requests/sec: 87522.6801

  Total data:   38873526 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.010 [999474]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.019 [35]    |
  0.028 [153]   |
  0.038 [0]     |
  0.047 [27]    |
  0.057 [48]    |
  0.066 [0]     |
  0.076 [52]    |
  0.085 [0]     |
  0.095 [65]    |


Latency distribution:
  10% in 0.0019 secs
  25% in 0.0023 secs
  50% in 0.0028 secs
  75% in 0.0033 secs
  90% in 0.0039 secs
  95% in 0.0044 secs
  99% in 0.0055 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.0945 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0442 secs
  resp wait:    0.0028 secs, 0.0001 secs, 0.0685 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0672 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

8. 2 core apiproxy, 4 worker 2 core, single nats

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        14.7657 secs
  Slowest:      0.0879 secs
  Fastest:      0.0002 secs
  Average:      0.0037 secs
  Requests/sec: 67714.5851

  Total data:   38873688 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.009 [999372]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.018 [125]   |
  0.027 [39]    |
  0.035 [221]   |
  0.044 [9]     |
  0.053 [19]    |
  0.062 [60]    |
  0.070 [0]     |
  0.079 [0]     |
  0.088 [9]     |


Latency distribution:
  10% in 0.0027 secs
  25% in 0.0032 secs
  50% in 0.0037 secs
  75% in 0.0042 secs
  90% in 0.0048 secs
  95% in 0.0053 secs
  99% in 0.0063 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0002 secs, 0.0879 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0558 secs
  resp wait:    0.0037 secs, 0.0002 secs, 0.0354 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0345 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

9. apiproxy, 4 worker 2 core, create multiple (8) connection on apiserver, single nats

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        11.8622 secs
  Slowest:      0.1487 secs
  Fastest:      0.0002 secs
  Average:      0.0030 secs
  Requests/sec: 84289.4330

  Total data:   38874336 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.015 [999398]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.030 [189]   |
  0.045 [80]    |
  0.060 [49]    |
  0.074 [82]    |
  0.089 [1]     |
  0.104 [17]    |
  0.119 [0]     |
  0.134 [14]    |
  0.149 [24]    |


Latency distribution:
  10% in 0.0017 secs
  25% in 0.0022 secs
  50% in 0.0028 secs
  75% in 0.0035 secs
  90% in 0.0044 secs
  95% in 0.0050 secs
  99% in 0.0067 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0002 secs, 0.1487 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0647 secs
  resp wait:    0.0025 secs, 0.0002 secs, 0.1475 secs
  resp read:    0.0003 secs, 0.0000 secs, 0.0362 secs

Status code distribution:
  [200] 999855 responses


###########################################################################

10. apiproxy, 4 worker 2 core, create multiple (8) connection on apiserver, cluster of 3 nats

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        11.2690 secs
  Slowest:      0.0964 secs
  Fastest:      0.0003 secs
  Average:      0.0029 secs
  Requests/sec: 88725.9153

  Total data:   38873613 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.010 [999261]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.020 [231]   |
  0.029 [36]    |
  0.039 [189]   |
  0.048 [20]    |
  0.058 [56]    |
  0.068 [32]    |
  0.077 [1]     |
  0.087 [7]     |
  0.096 [21]    |


Latency distribution:
  10% in 0.0017 secs
  25% in 0.0021 secs
  50% in 0.0027 secs
  75% in 0.0034 secs
  90% in 0.0042 secs
  95% in 0.0048 secs
  99% in 0.0063 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0003 secs, 0.0964 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0567 secs
  resp wait:    0.0024 secs, 0.0002 secs, 0.0550 secs
  resp read:    0.0003 secs, 0.0000 secs, 0.0325 secs

Status code distribution:
  [200] 999855 responses


###########################################################################

11. apiproxy, 1 worker 2 core, create multiple (8) connection on apiserver, cluster of 3 nats
hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        8.2034 secs
  Slowest:      0.0561 secs
  Fastest:      0.0002 secs
  Average:      0.0021 secs
  Requests/sec: 121883.4324

  Total data:   38874202 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.006 [960179]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.011 [38710] |■■
  0.017 [743]   |
  0.023 [35]    |
  0.028 [30]    |
  0.034 [27]    |
  0.039 [0]     |
  0.045 [0]     |
  0.051 [17]    |
  0.056 [113]   |


Latency distribution:
  10% in 0.0006 secs
  25% in 0.0008 secs
  50% in 0.0013 secs
  75% in 0.0029 secs
  90% in 0.0043 secs
  95% in 0.0054 secs
  99% in 0.0078 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0002 secs, 0.0561 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0279 secs
  resp wait:    0.0009 secs, 0.0001 secs, 0.0281 secs
  resp read:    0.0007 secs, 0.0000 secs, 0.0497 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

12. apiproxy, 1 callback worker, single nats

Summary:
  Total:        7.9362 secs
  Slowest:      0.0950 secs
  Fastest:      0.0002 secs
  Average:      0.0020 secs
  Requests/sec: 125986.4685

  Total data:   38873520 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.010 [997379]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.019 [2174]  |
  0.029 [105]   |
  0.038 [56]    |
  0.048 [0]     |
  0.057 [47]    |
  0.067 [27]    |
  0.076 [0]     |
  0.085 [25]    |
  0.095 [41]    |


Latency distribution:
  10% in 0.0006 secs
  25% in 0.0007 secs
  50% in 0.0012 secs
  75% in 0.0028 secs
  90% in 0.0043 secs
  95% in 0.0053 secs
  99% in 0.0077 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0002 secs, 0.0950 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0848 secs
  resp wait:    0.0008 secs, 0.0001 secs, 0.0527 secs
  resp read:    0.0007 secs, 0.0000 secs, 0.0856 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

13. apiproxy, 4 callback worker, single nats

Summary:
  Total:        10.9504 secs
  Slowest:      0.1029 secs
  Fastest:      0.0003 secs
  Average:      0.0028 secs
  Requests/sec: 91307.8431

  Total data:   38873665 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.011 [999023]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.021 [441]   |
  0.031 [97]    |
  0.041 [0]     |
  0.052 [223]   |
  0.062 [44]    |
  0.072 [5]     |
  0.082 [15]    |
  0.093 [0]     |
  0.103 [6]     |


Latency distribution:
  10% in 0.0016 secs
  25% in 0.0020 secs
  50% in 0.0025 secs
  75% in 0.0033 secs
  90% in 0.0042 secs
  95% in 0.0049 secs
  99% in 0.0068 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0003 secs, 0.1029 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0481 secs
  resp wait:    0.0022 secs, 0.0002 secs, 0.0535 secs
  resp read:    0.0004 secs, 0.0000 secs, 0.0487 secs

Status code distribution:
  [200] 999855 responses


###########################################################################

13. apiproxy, 2 core 1 callback worker, single nats

hey -n 1000000 -c 255 http://127.0.0.1:3000

Summary:
  Total:        8.0867 secs
  Slowest:      0.1232 secs
  Fastest:      0.0002 secs
  Average:      0.0020 secs
  Requests/sec: 123642.2649

  Total data:   38873987 bytes
  Size/request: 38 bytes

Response time histogram:
  0.000 [1]     |
  0.012 [999208]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.025 [350]   |
  0.037 [86]    |
  0.049 [95]    |
  0.062 [0]     |
  0.074 [0]     |
  0.086 [72]    |
  0.099 [11]    |
  0.111 [0]     |
  0.123 [32]    |


Latency distribution:
  10% in 0.0006 secs
  25% in 0.0007 secs
  50% in 0.0012 secs
  75% in 0.0029 secs
  90% in 0.0043 secs
  95% in 0.0054 secs
  99% in 0.0078 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0002 secs, 0.1232 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0843 secs
  resp wait:    0.0008 secs, 0.0001 secs, 0.0850 secs
  resp read:    0.0007 secs, 0.0000 secs, 0.1199 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

13. apiproxy, 2 core 4 callback worker, single nats

timedout a lot

*/
