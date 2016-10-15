# go-pubssed

![](images/pubssed-wof.png)

Listen to a Redis PubSub channel and then rebroadcast it over Server-Sent Events (SSE).

## Packages

### broker

```
import (
	"github.com/whosonfirst/go-pubssed/broker"
	"net/http"
)

br, _ := broker.NewBroker()
br.Start("localhost", 6379, "pubssed")

http_handler, err := br.HandlerFunc()

mux := http.NewServeMux()
mux.HandleFunc("/", http_handler)
http.ListenAndServe("localhost:8080", mux)
```

_Note that all error handling has been removed for the sake of brevity._

### listener

```
import (
	"github.com/whosonfirst/go-pubssed/listener"
	"log"
)

callback := func(msg string) error {
	log.Println(msg)
	return nil
}

lstnr, _ := listener.NewListener("http://localhost:8080", callback)
lstnr.Start()
```

_Note that all error handling has been removed for the sake of brevity._

## Tools

### pubssed-broadcast.go

```
./bin/pubssed-broadcast -h
Usage of ./bin/pubssed-broadcast:
  -clients int
    	Number of concurrent clients (default 200)
  -redis-channel string
    	Redis channel (default "pubssed")
  -redis-host string
    	Redis host (default "localhost")
  -redis-port int
    	Redis port (default 6379)
```
	
### pubssed-client.go

```
./bin/pubssed-client -h
Usage of ./bin/pubssed-client:
  -append-root string
    	The destination to write log files if the 'append' callback is invoked (default ".")
  -callback string
    	The callback to invoke when a SSE event is received (default "debug")
  -endpoint string
    	The pubssed endpoint you are connecting to
```

#### callbacks

##### append

Appends each SSE event a file named `YYYY/MM/DD/YYYYMMDDHH.txt` (where datetime specifics are determined based on the time the event is received).

##### debug

Logs each SSE event to STDOUT.

### pubssed-server.go

```
./bin/pubssed-server -h
Usage of ./bin/pubssed-server:
  -redis-channel string
    	Redis channel (default "pubssed")
  -redis-host string
    	Redis host (default "localhost")
  -redis-port int
    	Redis port (default 6379)
  -sse-endpoint string
    	SSE endpoint (default "/sse")
  -sse-host string
    	SSE host (default "localhost")
  -sse-port int
    	SSE port (default 8080)
```

## See also

* https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events
* https://github.com/whosonfirst/go-whosonfirst-webhookd
* https://github.com/whosonfirst/go-pubsocketd
