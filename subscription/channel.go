package subscription

import (
	"context"
)

type ChannelSubscription struct {
	Subscription
	channel chan string
}

func NewChannelSubscriptionWithChannel(ctx context.Context, ch chan string) (Subscription, error) {

	sub := &ChannelSubscription{
		channel: ch,
	}

	return sub, nil
}

func (sub *ChannelSubscription) Start(ctx context.Context, messages_ch chan string) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-sub.channel:
			messages_ch <- msg
		}
	}

	return nil
}

func (sub *ChannelSubscription) Close() error {
	return nil
}
