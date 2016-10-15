package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/whosonfirst/go-httpony/stats"
	"github.com/whosonfirst/go-pubssed/broker"
	"log"
	"net/http"
	"os"
)

func main() {

	var sse_host = flag.String("sse-host", "localhost", "SSE host")
	var sse_port = flag.Int("sse-port", 8080, "SSE port")

	var redis_host = flag.String("redis-host", "localhost", "Redis host")
	var redis_port = flag.Int("redis-port", 6379, "Redis port")
	var redis_channel = flag.String("redis-channel", "pubssed", "Redis channel")

	flag.Parse()

	br, err := broker.NewBroker()

	if err != nil {
		log.Fatal(err)
	}

	br.Start(*redis_host, *redis_port, *redis_channel)

	mux := http.NewServeMux()

	sse_handler, err := br.HandlerFunc()

	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/", sse_handler)

	stats_handler, err := stats.HandlerFunc(*sse_host)

	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/debug/vars", stats_handler)

	sse_endpoint := fmt.Sprintf("%s:%d", *sse_host, *sse_port)
	err = gracehttp.Serve(&http.Server{Addr: sse_endpoint, Handler: mux})

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
