package broker

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sfomuseum/go-pubsub/subscriber"
)

type Broker struct {
	clients      map[chan string]bool
	messages     chan string
	new_clients  chan chan string
	bunk_clients chan chan string
}

func NewBroker() (*Broker, error) {

	b := Broker{
		clients:      make(map[chan string]bool),
		messages:     make(chan string),
		new_clients:  make(chan (chan string)),
		bunk_clients: make(chan (chan string)),
	}

	return &b, nil
}

func (b *Broker) Start(ctx context.Context, sub subscriber.Subscriber) error {

	// set up the SSE monitor

	logger := slog.Default()
	logger.Debug("Start broker")

	go func() {

		for {

			select {

			case <-ctx.Done():
				logger.Debug("Broker received done signal, exiting")
				return

			case s := <-b.new_clients:

				b.clients[s] = true

			case s := <-b.bunk_clients:

				delete(b.clients, s)
				// We used to explicitly close(s) here. We should
				// really have been closing it below in the req.Context().Done()
				// block. Either way, don't since it is unnecessary and
				// seems to make Go start eating 100% of CPU...
				// log.Println("Removed client")

			case msg := <-b.messages:

				logger.Debug("Broadcast message to clients", "count", len(b.clients))

				for s, _ := range b.clients {
					s <- msg
				}
			}
		}
	}()

	// set up the PubSub monitor

	go func() {

		// something something error handling...
		logger.Debug("Listen")
		sub.Listen(ctx, b.messages)
	}()

	return nil
}

func (b *Broker) HandlerFunc() (http.HandlerFunc, error) {
	return b.HandlerFuncWithTimeout(nil)
}

func (b *Broker) HandlerFuncWithTimeout(ttl *time.Duration) (http.HandlerFunc, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		logger := slog.Default()
		logger = logger.With("remote addr", r.RemoteAddr)

		if ttl != nil {
			logger = logger.With("ttl", ttl)
		}

		t1 := time.Now()
		logger.Debug("Start broker HTTP handler", "time", t1)

		defer func() {
			logger.Debug("Finish handler", "time", time.Since(t1))
		}()

		fl, ok := w.(http.Flusher)

		if !ok {
			logger.Error("Writer does not support streaming")
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		messageChan := make(chan string)

		b.new_clients <- messageChan

		ctx := r.Context()

		if ttl != nil {

			c, cancel := context.WithTimeout(ctx, *ttl)
			defer cancel()

			ctx = c
		}

		notify := ctx.Done()

		go func() {
			<-notify
			b.bunk_clients <- messageChan
			// Don't close(messageChan) since it's unnecessary and
			// seems to cause CPU to spike to 100% Computers, amirite?
		}()

		// Set the headers related to event streaming.

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// https://stackoverflow.com/questions/27898622/server-sent-events-stopped-work-after-enabling-ssl-on-proxy
		// https://www.nginx.com/resources/wiki/start/topics/examples/x-accel/#X-Accel-Buffering

		w.Header().Set("X-Accel-Buffering", "no")

		// For CORS stuff please use https://github.com/rs/cors

		for {

			select {
			case <-notify:
				logger.Debug("Handler received signal, stopping handler")
				return
			case msg := <-messageChan:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				fl.Flush()
			}
		}

	}

	return http.HandlerFunc(f), nil
}
