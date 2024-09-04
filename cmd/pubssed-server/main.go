package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/sfomuseum/go-pubsub/subscriber"
	"github.com/whosonfirst/go-pubssed/broker"
)

//go:embed index.html index.js
var FS embed.FS

func main() {

	var sse_host = flag.String("sse-host", "localhost", "SSE host")
	var sse_port = flag.Int("sse-port", 8080, "SSE port")
	var sse_endpoint = flag.String("sse-endpoint", "/sse", "SSE endpoint")

	var subscription_uri = flag.String("subscription-uri", "redis://?host=localhost&port=6379&channel=pubssed", "...")

	var verbose = flag.Bool("verbose", false, "Enable verbose (debug) logging")
	enable_demo := flag.Bool("enable-demo", false, "...")

	flag.Parse()

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	slog.Debug("SERVER", "sub", *subscription_uri)
	sub, err := subscriber.NewSubscriber(ctx, *subscription_uri)

	if err != nil {
		log.Fatalf("Failed to create subscription for '%s', %v", *subscription_uri, err)
	}

	br, err := broker.NewBroker()

	if err != nil {
		log.Fatal(err)
	}

	br.Start(ctx, sub)

	mux := http.NewServeMux()

	sse_handler, err := br.HandlerFunc()

	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc(*sse_endpoint, sse_handler)

	if *enable_demo {
		http_fs := http.FS(FS)
		fs_handler := http.FileServer(http_fs)
		mux.Handle("/", fs_handler)
	}

	sse_addr := fmt.Sprintf("%s:%d", *sse_host, *sse_port)
	log.Printf("Listening on %s\n", sse_addr)

	err = http.ListenAndServe(sse_addr, mux)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
