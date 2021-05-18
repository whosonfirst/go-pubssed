package broker

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "log"
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

func (b *Broker) Start(redis_host string, redis_port int, redis_channel string) error {

	redis_endpoint := fmt.Sprintf("%s:%d", redis_host, redis_port)

	/*

		we're going to set up the redis/pubsub clients again below inside the scope
		of a Go routine; we're doing it here to perform basic sanity checking on the
		subscription (20161014/thisisaaronland)

	*/

	redis_client := redis.NewClient(&redis.Options{
		Addr: redis_endpoint,
	})

	defer redis_client.Close()

	pubsub_client := redis_client.PSubscribe(redis_channel)
	defer pubsub_client.Close()

	// set up the SSE monitor

	go func() {

		for {

			select {

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

		redis_client := redis.NewClient(&redis.Options{
			Addr: redis_endpoint,
		})

		defer redis_client.Close()

		pubsub_client := redis_client.PSubscribe(redis_channel)
		defer pubsub_client.Close()

		for {

			i, _ := pubsub_client.Receive()

			// log.Println("received message", i)

			if msg, _ := i.(*redis.Message); msg != nil {

				// log.Println("relay message", msg.Payload)
				b.messages <- msg.Payload
			}
		}
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

		notify := w.(http.CloseNotifier).CloseNotify()

		go func() {
			<-notify
			b.bunk_clients <- messageChan
			// log.Println("HTTP connection just closed.")
		}()

		// Set the headers related to event streaming.

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Don't close the connection, instead loop 10 times,
		// sending messages and flushing the response each time
		// there is a new message to send along.
		//
		// NOTE: we could loop endlessly; however, then you
		// could not easily detect clients that dettach and the
		// server would continue to send them messages long after
		// they're gone due to the "keep-alive" header.  One of
		// the nifty aspects of SSE is that clients automatically
		// reconnect when they lose their connection.
		//
		// A better way to do this is to use the CloseNotifier
		// interface that will appear in future releases of
		// Go (this is written as of 1.0.3):
		// https://code.google.com/p/go/source/detail?name=3292433291b2

		for {

			msg, open := <-messageChan

			if !open {
				break
			}

			fmt.Fprintf(w, "data: %s\n\n", msg)
			fl.Flush()
		}

		// log.Println("Finished HTTP request at ", r.URL.Path)
	}

	return http.HandlerFunc(f), nil
}
