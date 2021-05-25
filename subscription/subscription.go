package subscription

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
)

type Subscription interface {
	Start(context.Context, chan string) error
	Close() error
}

type SubscriptionInitializeFunc func(ctx context.Context, uri string) (Subscription, error)

var subscriptions roster.Roster

func ensureSpatialRoster() error {

	if subscriptions == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		subscriptions = r
	}

	return nil
}

func RegisterSubscription(ctx context.Context, scheme string, f SubscriptionInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return subscriptions.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range subscriptions.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewSubscription(ctx context.Context, uri string) (Subscription, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := subscriptions.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(SubscriptionInitializeFunc)
	return f(ctx, uri)
}
