package subscriber

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-redis/redis/v8"
)

type RedisSubscriber struct {
	Subscriber
	redis_client  *redis.Client
	redis_channel string
}

func init() {
	ctx := context.Background()
	RegisterSubscriber(ctx, "redis", NewRedisSubscriber)
}

func NewRedisSubscriber(ctx context.Context, uri string) (Subscriber, error) {

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

	s := &RedisSubscriber{
		redis_client:  redis_client,
		redis_channel: channel,
	}

	return s, nil
}

func (s *RedisSubscriber) Listen(ctx context.Context, messages_ch chan string) error {

	pubsub_client := s.redis_client.PSubscribe(ctx, s.redis_channel)
	defer pubsub_client.Close()

	for {

		i, err := pubsub_client.Receive(ctx)

		if err != nil {
			return fmt.Errorf("Failed to receive message, %w", err)
		}

		if msg, _ := i.(*redis.Message); msg != nil {
			messages_ch <- msg.Payload
		}
	}

	return nil
}

func (s *RedisSubscriber) Close() error {
	return s.redis_client.Close()
}
