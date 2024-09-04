package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {

	var clients = flag.Int("clients", 200, "Number of concurrent clients")

	var redis_host = flag.String("redis-host", "localhost", "Redis host")
	var redis_port = flag.Int("redis-port", 6379, "Redis port")
	var redis_channel = flag.String("redis-channel", "pubssed", "Redis channel")

	var verbose = flag.Bool("verbose", false, "Enable verbose (debug) logging")

	flag.Parse()

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	redis_endpoint := fmt.Sprintf("%s:%d", *redis_host, *redis_port)

	redis_client := redis.NewClient(&redis.Options{
		Addr: redis_endpoint,
	})

	defer redis_client.Close()

	ch := make(chan bool, *clients)

	for i := 0; i < *clients; i++ {
		ch <- true
	}

	logger := slog.Default()
	logger = logger.With("channel", *redis_channel)

	for {

		<-ch

		now := fmt.Sprintf("%v", time.Now())
		logger.Info("Publish", "message", now)

		redis_client.Publish(ctx, *redis_channel, now)

		ch <- true
	}

	os.Exit(0)
}
