// Package server provides a sample HTTP/Websocket server which registers
// itself in consul using one or more url prefixes to demonstrate and
// test the automatic fabio routing table update.
//
// During startup the server performs the following steps:
//
// * Add a handler for each prefix which provides a unique
//   response for that instance and endpoint
// * Add a `/health` handler for the consul health check
// * Register the service in consul with the listen address,
//   a health check under the given name and with one `urlprefix-`
//   tag per prefix
// * Install a signal handler to deregister the service on exit
//
// If the protocol is set to "ws" the registered endpoints function
// as websocket echo servers.
//
// Example:
//
//   # http server
//   ./server -addr 127.0.0.1:5000 -name svc-a -prefix /foo -prefix /bar
//   ./server -addr 127.0.0.1:5001 -name svc-b -prefix /baz -prefix /bar
//   ./server -addr 127.0.0.1:5002 -name svc-c -prefix "/gogl redirect=301,https://www.google.de/"
//
//   # https server
//   ./server -addr 127.0.0.1:5000 -name svc-a -proto https -certFile ... -keyFile ... -prefix /foo
//   ./server -addr 127.0.0.1:5000 -name svc-a -proto https -certFile ... -keyFile ... -prefix "/foo tlsskipverify=true"
//
//   # websocket server
//   ./server -addr 127.0.0.1:6000 -name ws-a -proto ws -prefix /echo1 -prefix /echo2
//
//   # tcp server
//   ./server -addr 127.0.0.1:7000 -name tcp-a -proto tcp -prefix :1234
//
// source: https://github.com/fabiolb/fabio/blob/master/demo/server/server.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fabiolb/fabio/proxy/tcp"
	"github.com/hashicorp/consul/api"
	"golang.org/x/net/websocket"
)

type Args struct {
	addr     string
	consul   string
	name     string
	proto    string
	token    string
	certFile string
	keyFile  string
	status   int
	prefixes []string
	tags     []string
}

type stringsVar []string

func (s *stringsVar) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s stringsVar) String() string {
	return strings.Join(s, " ")
}

// setup:
// curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
// sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
// sudo apt install consul
// go install github.com/fabiolb/fabio@latest

// start:
// sudo consul agent -dev --data-dir=/tmp/consul
// fabio
// go run main.go -addr 172.17.0.1:5000 -name svc-a -prefix /foo -consul 127.0.0.1:8500

// http://127.0.0.1:8500 for consul web GUI/REST API
// http://127.0.0.1:9998 for fabio web GUI
// http://127.0.0.1:9999 for public facing HTTP

func main() {
	var args Args

	flag.StringVar(&args.addr, "addr", "127.0.0.1:5000", "host:port of the service")
	flag.StringVar(&args.consul, "consul", "127.0.0.1:8500", "host:port of the consul agent")
	flag.StringVar(&args.name, "name", filepath.Base(os.Args[0]), "name of the service")
	flag.StringVar(&args.proto, "proto", "http", "protocol for endpoints: http, ws or tcp")
	flag.StringVar(&args.token, "token", "", "consul ACL token")
	flag.StringVar(&args.certFile, "cert", "", "path to cert file")
	flag.StringVar(&args.keyFile, "key", "", "path to key file")
	flag.IntVar(&args.status, "status", http.StatusOK, "http status code")
	flag.Var((*stringsVar)(&args.prefixes), "prefix", "'host/path' or ':port' prefix to register. Can be specified multiple times")
	flag.Var((*stringsVar)(&args.tags), "tags", "additional tags to register in consul. Can be specified multiple times")
	flag.Parse()

	if len(args.prefixes) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if (args.proto == "https" || args.proto == "wss") && args.certFile == "" {
		log.Fatalf("Proto %s requires a certificate. Please provide -cert/-key", args.proto)
	}

	type server interface {
		ListenAndServe() error
		ListenAndServeTLS(certFile, keyFile string) error
	}

	var srv server
	var tags []string
	var check *api.AgentServiceCheck
	switch args.proto {
	case "http", "https", "ws", "wss":
		srv, tags, check = newHTTPServer(args)
	case "tcp":
		srv, tags, check = newTCPServer(args)
	default:
		log.Fatal("Invalid protocol ", args.proto)
	}

	// start server
	go func() {
		var err error
		if args.certFile != "" {
			err = srv.ListenAndServeTLS(args.certFile, args.keyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil {
			log.Fatal(err)
		}
	}()

	// get host and port as string/int
	host, portstr, err := net.SplitHostPort(args.addr)
	if err != nil {
		log.Fatal(err)
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		log.Fatal(err)
	}

	// register service with health check
	serviceID := args.name + "-" + args.addr
	service := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    args.name,
		Port:    port,
		Address: host,
		Tags:    tags,
		Check:   check,
	}

	config := &api.Config{Address: args.consul, Scheme: "http", Token: args.token}
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Agent().ServiceRegister(service); err != nil {
		log.Fatal(err)
	}
	log.Printf("Registered %s service %q in consul with tags %q", args.proto, args.name, strings.Join(tags, ","))

	// run until we get a signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	// deregister service
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		log.Fatal(err)
	}
	log.Printf("Deregistered service %q in consul", args.name)
}

func newHTTPServer(args Args) (*http.Server, []string, *api.AgentServiceCheck) {
	addr, proto, tags := args.addr, args.proto, args.tags

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", 404)
		log.Printf("%s -> 404", r.URL)
	})

	for _, p := range args.prefixes {
		uri := strings.Fields(p)[0]
		switch proto {
		case "http", "https":
			mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				w.WriteHeader(args.status)
				fmt.Fprintf(w, "Serving %s via %s from %s on %s\n", r.RequestURI, scheme, args.name, addr)
			})
		case "ws", "wss":
			mux.Handle(uri, websocket.Handler(WSEchoServer))
		}

		tag := "urlprefix-" + p
		if proto == "https" || proto == "wss" {
			tag += " proto=https"
		}
		tags = append(tags, tag)
	}

	checkScheme := "http"
	if args.certFile != "" {
		checkScheme = "https"
	}
	check := &api.AgentServiceCheck{
		HTTP:          checkScheme + "://" + addr + "/health",
		Interval:      "1s",
		Timeout:       "1s",
		TLSSkipVerify: true,
	}
	return &http.Server{Addr: addr, Handler: mux}, tags, check
}

func WSEchoServer(ws *websocket.Conn) {
	addr := ws.LocalAddr().String()
	pfx := []byte("[" + addr + "] ")

	log.Printf("ws connect on %s", addr)

	// the following could be done with io.Copy(ws, ws)
	// but I want to add some meta data
	var msg = make([]byte, 1024)
	for {
		n, err := ws.Read(msg)
		if err != nil && err != io.EOF {
			log.Printf("ws error on %s. %s", addr, err)
			break
		}
		_, err = ws.Write(append(pfx, msg[:n]...))
		if err != nil && err != io.EOF {
			log.Printf("ws error on %s. %s", addr, err)
			break
		}
	}
	log.Printf("ws disconnect on %s", addr)
}

func newTCPServer(args Args) (*tcp.Server, []string, *api.AgentServiceCheck) {
	tags := args.tags
	for _, p := range args.prefixes {
		tags = append(tags, "urlprefix-"+p+" proto=tcp")
	}
	check := &api.AgentServiceCheck{TCP: args.addr, Interval: "2s", Timeout: "1s"}
	return &tcp.Server{Addr: args.addr, Handler: tcp.HandlerFunc(TCPEchoHandler)}, tags, check
}

func TCPEchoHandler(c net.Conn) error {
	defer c.Close()

	addr := c.LocalAddr().String()
	_, err := fmt.Fprintf(c, "[%s] Welcome\n", addr)
	if err != nil {
		return err
	}

	for {
		line, _, err := bufio.NewReader(c).ReadLine()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(c, "[%s] %s\n", addr, string(line))
		if err != nil {
			return err
		}
	}
}

// https://devopscube.com/setup-consul-cluster-guide/

/*
benchmark scenario:

###########################################################################

1. direct handling health without fabio

hey -n 1000000 -c 255 http://172.17.0.1:5000/health

Summary:
  Total:        5.1817 secs
  Slowest:      0.0909 secs
  Fastest:      0.0000 secs
  Average:      0.0013 secs
  Requests/sec: 192958.3425

  Total data:   2999565 bytes
  Size/request: 3 bytes

Response time histogram:
  0.000 [1]     |
  0.009 [996345]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.018 [3209]  |
  0.027 [39]    |
  0.036 [149]   |
  0.045 [11]    |
  0.055 [43]    |
  0.064 [18]    |
  0.073 [0]     |
  0.082 [2]     |
  0.091 [38]    |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0001 secs
  75% in 0.0021 secs
  90% in 0.0039 secs
  95% in 0.0046 secs
  99% in 0.0075 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0909 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0538 secs
  resp wait:    0.0001 secs, 0.0000 secs, 0.0823 secs
  resp read:    0.0006 secs, 0.0000 secs, 0.0541 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

2. direct handling foo without fabio

hey -n 1000000 -c 255 http://172.17.0.1:5000/foo

Summary:
  Total:        5.0742 secs
  Slowest:      0.0644 secs
  Fastest:      0.0000 secs
  Average:      0.0013 secs
  Requests/sec: 197047.9124

  Total data:   51992460 bytes
  Size/request: 52 bytes

Response time histogram:
  0.000 [1]     |
  0.006 [982587]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.013 [16855] |■
  0.019 [257]   |
  0.026 [71]    |
  0.032 [19]    |
  0.039 [39]    |
  0.045 [0]     |
  0.052 [0]     |
  0.058 [0]     |
  0.064 [26]    |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0001 secs
  75% in 0.0021 secs
  90% in 0.0038 secs
  95% in 0.0046 secs
  99% in 0.0074 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0644 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0364 secs
  resp wait:    0.0001 secs, 0.0000 secs, 0.0374 secs
  resp read:    0.0006 secs, 0.0000 secs, 0.0364 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

3. proxied /foo with fabio

hey -n 1000000 -c 255 http://127.0.0.1:9999/foo

Summary:
  Total:        15.2035 secs
  Slowest:      0.1030 secs
  Fastest:      0.0001 secs
  Average:      0.0038 secs
  Requests/sec: 65764.9021

  Total data:   51992460 bytes
  Size/request: 52 bytes

Response time histogram:
  0.000 [1]     |
  0.010 [949783]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.021 [48836] |■■
  0.031 [753]   |
  0.041 [111]   |
  0.052 [241]   |
  0.062 [6]     |
  0.072 [29]    |
  0.082 [0]     |
  0.093 [0]     |
  0.103 [95]    |


Latency distribution:
  10% in 0.0003 secs
  25% in 0.0006 secs
  50% in 0.0036 secs
  75% in 0.0056 secs
  90% in 0.0085 secs
  95% in 0.0104 secs
  99% in 0.0145 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.1030 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0707 secs
  resp wait:    0.0037 secs, 0.0001 secs, 0.0703 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0701 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

4. proxied /foo with fabio, 2 instance

hey -n 1000000 -c 255 http://127.0.0.1:9999/foo

Summary:
  Total:        14.6571 secs
  Slowest:      0.1320 secs
  Fastest:      0.0001 secs
  Average:      0.0037 secs
  Requests/sec: 68216.5154

  Total data:   51992460 bytes
  Size/request: 52 bytes

Response time histogram:
  0.000 [1]     |
  0.013 [985686]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.027 [13702] |■
  0.040 [232]   |
  0.053 [176]   |
  0.066 [31]    |
  0.079 [1]     |
  0.092 [19]    |
  0.106 [0]     |
  0.119 [0]     |
  0.132 [7]     |


Latency distribution:
  10% in 0.0003 secs
  25% in 0.0006 secs
  50% in 0.0030 secs
  75% in 0.0057 secs
  90% in 0.0080 secs
  95% in 0.0101 secs
  99% in 0.0141 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.1320 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0805 secs
  resp wait:    0.0036 secs, 0.0001 secs, 0.0805 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0439 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

5. proxied /foo with fabio, 1 instance 2 core

hey -n 1000000 -c 255 http://127.0.0.1:9999/foo

Summary:
  Total:        17.7469 secs
  Slowest:      0.1719 secs
  Fastest:      0.0001 secs
  Average:      0.0045 secs
  Requests/sec: 56339.5518

  Total data:   51992460 bytes
  Size/request: 52 bytes

Response time histogram:
  0.000 [1]     |
  0.017 [998327]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.034 [921]   |
  0.052 [380]   |
  0.069 [121]   |
  0.086 [39]    |
  0.103 [55]    |
  0.120 [1]     |
  0.138 [7]     |
  0.155 [0]     |
  0.172 [3]     |


Latency distribution:
  10% in 0.0011 secs
  25% in 0.0024 secs
  50% in 0.0044 secs
  75% in 0.0060 secs
  90% in 0.0077 secs
  95% in 0.0092 secs
  99% in 0.0123 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.1719 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0669 secs
  resp wait:    0.0044 secs, 0.0001 secs, 0.1037 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0661 secs

Status code distribution:
  [200] 999855 responses


###########################################################################

6. proxied /foo with fabio, 2 instance 2 core

hey -n 1000000 -c 255 http://127.0.0.1:9999/foo

Summary:
  Total:        16.5822 secs
  Slowest:      0.1153 secs
  Fastest:      0.0001 secs
  Average:      0.0042 secs
  Requests/sec: 60296.9714

  Total data:   51992460 bytes
  Size/request: 52 bytes

Response time histogram:
  0.000 [1]     |
  0.012 [966182]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.023 [32700] |■
  0.035 [508]   |
  0.046 [115]   |
  0.058 [18]    |
  0.069 [246]   |
  0.081 [75]    |
  0.092 [0]     |
  0.104 [0]     |
  0.115 [10]    |


Latency distribution:
  10% in 0.0005 secs
  25% in 0.0010 secs
  50% in 0.0035 secs
  75% in 0.0064 secs
  90% in 0.0086 secs
  95% in 0.0104 secs
  99% in 0.0152 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.1153 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0776 secs
  resp wait:    0.0041 secs, 0.0001 secs, 0.0772 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0577 secs

Status code distribution:
  [200] 999855 responses


###########################################################################

7. proxied /foo with fabio 8 core, 1 instance 2 core

Summary:
  Total:        16.6727 secs
  Slowest:      0.1435 secs
  Fastest:      0.0001 secs
  Average:      0.0042 secs
  Requests/sec: 59969.5206

  Total data:   51992460 bytes
  Size/request: 52 bytes

Response time histogram:
  0.000 [1]     |
  0.014 [994778]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.029 [4589]  |
  0.043 [148]   |
  0.057 [127]   |
  0.072 [87]    |
  0.086 [63]    |
  0.100 [50]    |
  0.115 [7]     |
  0.129 [0]     |
  0.144 [5]     |


Latency distribution:
  10% in 0.0011 secs
  25% in 0.0021 secs
  50% in 0.0039 secs
  75% in 0.0055 secs
  90% in 0.0075 secs
  95% in 0.0092 secs
  99% in 0.0126 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.1435 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0324 secs
  resp wait:    0.0042 secs, 0.0001 secs, 0.0914 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0303 secs

Status code distribution:
  [200] 999855 responses



###########################################################################

8. proxied /foo with fabio 8 core, 2 instance 2 core

Summary:
  Total:        16.0828 secs
  Slowest:      0.1045 secs
  Fastest:      0.0001 secs
  Average:      0.0041 secs
  Requests/sec: 62169.2256

  Total data:   51992460 bytes
  Size/request: 52 bytes

Response time histogram:
  0.000 [1]     |
  0.011 [971653]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.021 [27260] |■
  0.031 [400]   |
  0.042 [279]   |
  0.052 [37]    |
  0.063 [96]    |
  0.073 [78]    |
  0.084 [21]    |
  0.094 [0]     |
  0.105 [30]    |


Latency distribution:
  10% in 0.0009 secs
  25% in 0.0018 secs
  50% in 0.0037 secs
  75% in 0.0056 secs
  90% in 0.0076 secs
  95% in 0.0092 secs
  99% in 0.0126 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.1045 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0393 secs
  resp wait:    0.0040 secs, 0.0001 secs, 0.0790 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0595 secs

Status code distribution:
  [200] 999855 responses

###########################################################################

9. proxied /foo with fabio 8 core, 4 instance 2 core

hey -n 1000000 -c 255 http://127.0.0.1:9999/foo

Summary:
  Total:        15.4528 secs
  Slowest:      0.0743 secs
  Fastest:      0.0001 secs
  Average:      0.0039 secs
  Requests/sec: 64703.8253

  Total data:   51992460 bytes
  Size/request: 52 bytes

Response time histogram:
  0.000 [1]     |
  0.008 [912009]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.015 [85624] |■■■■
  0.022 [1863]  |
  0.030 [225]   |
  0.037 [31]    |
  0.045 [56]    |
  0.052 [43]    |
  0.059 [0]     |
  0.067 [0]     |
  0.074 [3]     |


Latency distribution:
  10% in 0.0009 secs
  25% in 0.0017 secs
  50% in 0.0035 secs
  75% in 0.0055 secs
  90% in 0.0073 secs
  95% in 0.0087 secs
  99% in 0.0122 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.0743 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0510 secs
  resp wait:    0.0038 secs, 0.0001 secs, 0.0490 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0508 secs

Status code distribution:
  [200] 999855 responses
*/
