package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/sfomuseum/go-pubsub/publisher"
)

func main() {

	var clients = flag.Int("clients", 200, "Number of concurrent clients")

	var publisher_uri = flag.String("publisher-uri", "redis://?host=localhost&port=6379&channel=pubssed", "...")

	var verbose = flag.Bool("verbose", false, "Enable verbose (debug) logging")

	flag.Parse()

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	pub, err := publisher.NewPublisher(ctx, *publisher_uri)

	if err != nil {
		log.Fatalf("Failed to create publisher, %w", err)
	}

	defer pub.Close()

	ch := make(chan bool, *clients)

	for i := 0; i < *clients; i++ {
		ch <- true
	}

	logger := slog.Default()
	logger = logger.With("publisher", *publisher_uri)

	for {

		<-ch

		now := fmt.Sprintf("%v", time.Now())
		logger.Info("Publish", "message", now)

		pub.Publish(ctx, now)

		ch <- true
	}

	os.Exit(0)
}
