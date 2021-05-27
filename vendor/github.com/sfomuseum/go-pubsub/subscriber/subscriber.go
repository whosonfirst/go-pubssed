package subscriber

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"net/url"
	"sort"
	"strings"
)

type Subscriber interface {
	Listen(context.Context, chan string) error
	Close() error
}

type SubscriberInitializeFunc func(ctx context.Context, uri string) (Subscriber, error)

var subscribers roster.Roster

func ensureSpatialRoster() error {

	if subscribers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		subscribers = r
	}

	return nil
}

func RegisterSubscriber(ctx context.Context, scheme string, f SubscriberInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return subscribers.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range subscribers.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewSubscriber(ctx context.Context, uri string) (Subscriber, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := subscribers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(SubscriberInitializeFunc)
	return f(ctx, uri)
}
