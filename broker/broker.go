package broker

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-pubsub/subscriber"
	"log"
	"net/http"
	"os"
)

type Broker struct {
	clients      map[chan string]bool
	messages     chan string
	new_clients  chan chan string
	bunk_clients chan chan string
	Logger       *log.Logger
}

func NewBroker() (*Broker, error) {

	logger := log.New(os.Stdout, "[pubssed] ", log.LstdFlags)

	b := Broker{
		clients:      make(map[chan string]bool),
		messages:     make(chan string),
		new_clients:  make(chan (chan string)),
		bunk_clients: make(chan (chan string)),
		Logger:       logger,
	}

	return &b, nil
}

func (b *Broker) Start(ctx context.Context, sub subscriber.Subscriber) error {

	// set up the SSE monitor

	go func() {

		for {

			select {

			case <-ctx.Done():
				// log.Println("Done")
				return

			case s := <-b.new_clients:

				b.clients[s] = true
				// log.Println("Added new client")

			case s := <-b.bunk_clients:

				delete(b.clients, s)
				// We used to explicitly close(s) here. We should
				// really have been closing it below in the req.Context().Done()
				// block. Either way, don't since it is unnecessary and
				// seems to make Go start eating 100% of CPU...
				// log.Println("Removed client")

			case msg := <-b.messages:

				for s, _ := range b.clients {
					s <- msg
				}

				// log.Printf("Broadcast message to %d clients", len(b.clients))
			}
		}
	}()

	// set up the PubSub monitor

	go func() {

		// something something error handling...

		sub.Listen(ctx, b.messages)
	}()

	return nil
}

func (b *Broker) HandlerFunc() (http.HandlerFunc, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		b.Logger.Println("SSE start handler")

		defer func() {
			b.Logger.Println("SSE finish handler")
		}()

		fl, ok := w.(http.Flusher)

		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		messageChan := make(chan string)

		b.new_clients <- messageChan

		notify := r.Context().Done()

		go func() {
			<-notify
			b.Logger.Println("SSE HTTP connection just closed.")
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
				b.Logger.Println("SSE stop handler")
				return
			case msg := <-messageChan:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				fl.Flush()
			}
		}

	}

	return http.HandlerFunc(f), nil
}
