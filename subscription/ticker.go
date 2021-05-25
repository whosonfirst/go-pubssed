package subscription

import (
	"context"
	"net/url"
	"time"
)

type TickerSubscription struct {
	Subscription
	ticker *time.Ticker
}

func init() {
	ctx := context.Background()
	RegisterSubscription(ctx, "ticker", NewTickerSubscription)
}

func NewTickerSubscription(ctx context.Context, uri string) (Subscription, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	t := time.NewTicker(1 * time.Second)

	sub := &TickerSubscription{
		ticker: t,
	}

	return sub, nil
}

func (sub *TickerSubscription) Start(ctx context.Context, messages_ch chan string) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-sub.ticker.C:
			messages_ch <- t.Format(time.RFC3339)
		}
	}

	return nil
}

func (sub *TickerSubscription) Close() error {
	sub.ticker.Stop()
	return nil
}
