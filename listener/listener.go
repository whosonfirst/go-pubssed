package listener

import (
	"bufio"
	"bytes"
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
		return err
	}

	req.Header.Set("Accept", "text/event-stream")

	res, err := l.client.Do(req)

	if err != nil {
		return err
	}

	br := bufio.NewReader(res.Body)
	defer res.Body.Close()

	delim := []byte{':', ' '}

	for {
		bs, err := br.ReadBytes('\n')

		if err != nil {

			if err.Error() == "EOF" {
			   	return l.Start()
			}
			
			return err
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

		msg := string(bytes.TrimSpace(spl[2]))

		err = l.callback(msg)

		if err != nil {
			return err
		}
	}

	return nil
}
