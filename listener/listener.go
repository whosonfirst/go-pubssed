package listener

import (
	"bufio"
	"bytes"
	"fmt"
	_ "log"
	"net/http"
)

type Listener struct {
	client   *http.Client
	endpoint string
	callback ListenerFunc
}

type ListenerFunc func(string) error

func NewListener(endpoint string, callback ListenerFunc) (*Listener, error) {

	client := &http.Client{}

	l := Listener{
		client:   client,
		endpoint: endpoint,
		callback: callback,
	}

	return &l, nil
}

func (l *Listener) Start() error {

	req, err := http.NewRequest("GET", l.endpoint, nil)

	if err != nil {
		return fmt.Errorf("Failed to create new request, %w", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")
	//req.Header.Set("Accept", "text/event-stream")

	res, err := l.client.Do(req)

	if err != nil {
		return fmt.Errorf("Failed to request %s, %w", l.endpoint, err)
	}

	br := bufio.NewReader(res.Body)
	defer res.Body.Close()

	delim := []byte{':', ' '}

	for {
		bs, err := br.ReadBytes('\n')

		if err != nil {
			return fmt.Errorf("Failed to read bytes, %w", err)
		}

		if len(bs) < 2 {
			continue
		}

		spl := bytes.Split(bs, delim)

		if len(spl) < 2 {
			continue
		}

		ctx := string(spl[0])

		if ctx != "data" {
			continue
		}

		msg := string(bytes.TrimSpace(spl[1]))

		err = l.callback(msg)

		if err != nil {
			return fmt.Errorf("Listener callback failed, %w", err)
		}
	}

	return nil
}
