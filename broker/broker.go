package broker

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-pubsub/subscriber"
	"log"
	"net/http"
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
				close(s)

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
			log.Println("HTTP connection just closed.")
			b.bunk_clients <- messageChan
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

			msg := <-messageChan

			fmt.Fprintf(w, "data: %s\n\n", msg)
			fl.Flush()
		}

		log.Println("Finished HTTP request at ", r.URL.Path)
	}

	return http.HandlerFunc(f), nil
}
