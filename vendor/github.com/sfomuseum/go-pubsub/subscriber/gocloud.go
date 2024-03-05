package subscriber

// https://gocloud.dev/howto/pubsub/subscribe/

import (
	"context"

	"gocloud.dev/pubsub"
)

type GoCloudSubscriber struct {
	Subscriber
	subscription *pubsub.Subscription
}

func init() {

	ctx := context.Background()
	RegisterGoCloudSubscribers(ctx)
}

func RegisterGoCloudSubscribers(ctx context.Context) error {

	for _, scheme := range pubsub.DefaultURLMux().SubscriptionSchemes() {

		err := RegisterSubscriber(ctx, scheme, NewGoCloudSubscriber)

		if err != nil {
			panic(err)
		}
	}

	return nil
}

func NewGoCloudSubscriber(ctx context.Context, uri string) (Subscriber, error) {

	subscription, err := pubsub.OpenSubscription(ctx, uri)

	if err != nil {
		return nil, err
	}

	sub := &GoCloudSubscriber{
		subscription: subscription,
	}

	return sub, err
}

func (sub *GoCloudSubscriber) Listen(ctx context.Context, msg_ch chan string) error {

	for {

		msg, err := sub.subscription.Receive(ctx)

		if err != nil {
			return err
		}

		go msg.Ack()

		msg_ch <- string(msg.Body)
	}

	return nil
}

func (sub *GoCloudSubscriber) Close() error {
	ctx := context.Background()
	return sub.subscription.Shutdown(ctx)
}
