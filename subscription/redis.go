package subscription

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/url"
)

type RedisSubscription struct {
	Subscription
	redis_client  *redis.Client
	redis_channel string
}

func init() {
	ctx := context.Background()
	RegisterSubscription(ctx, "redis", NewRedisSubscription)
}

func NewRedisSubscription(ctx context.Context, uri string) (Subscription, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	host := q.Get("host")
	port := q.Get("port")
	channel := q.Get("channel")

	addr := fmt.Sprintf("%s:%s", host, port)

	redis_client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	s := &RedisSubscription{
		redis_client:  redis_client,
		redis_channel: channel,
	}

	return s, nil
}

func (s *RedisSubscription) Start(ctx context.Context, messages_ch chan string) error {

	pubsub_client := s.redis_client.PSubscribe(ctx, s.redis_channel)
	defer pubsub_client.Close()

	for {

		i, _ := pubsub_client.Receive(ctx)
		// log.Println("received message", i)

		if msg, _ := i.(*redis.Message); msg != nil {
			log.Println("relay message", msg.Payload)
			messages_ch <- msg.Payload
		}
	}

	return nil
}

func (s *RedisSubscription) Close() error {
	return s.redis_client.Close()
}
