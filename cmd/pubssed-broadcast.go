package main

import (
	"flag"
	"fmt"
	"gopkg.in/redis.v1"
	_ "log"
	"os"
	"time"
)

func main() {

	var clients = flag.Int("clients", 200, "Number of concurrent clients")

	var redis_host = flag.String("redis-host", "localhost", "Redis host")
	var redis_port = flag.Int("redis-port", 6379, "Redis port")
	var redis_channel = flag.String("redis-channel", "pubssed", "Redis channel")

	flag.Parse()

	redis_endpoint := fmt.Sprintf("%s:%d", *redis_host, *redis_port)

	redis_client := redis.NewTCPClient(&redis.Options{
		Addr: redis_endpoint,
	})

	defer redis_client.Close()

	ch := make(chan bool, *clients)

	for i := 0; i < *clients; i++ {
		ch <- true
	}

	for {

		<-ch

		now := fmt.Sprintf("%v", time.Now())
		redis_client.Publish(*redis_channel, now)

		ch <- true
	}

	os.Exit(0)
}
