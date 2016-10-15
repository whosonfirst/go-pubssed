package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-pubssed/listener"
	"log"
	"os"
	"time"
)

func main() {

	var endpoint = flag.String("endpoint", "", "...")
	var callback = flag.String("callback", "debug", "...")

	flag.Parse()

	var f listener.ListenerFunc

	if *callback == "debug" {

		f = func(msg string) error {

			log.Println(msg)
			return nil
		}

	} else if *callback == "append" {

		f = func(msg string) error {

			now := time.Now()
			fname := fmt.Sprintf("%04d%02d%02d%02d.txt", now.Year(), now.Month(), now.Day(), now.Hour())

			fh, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND, 0644)

			if err != nil {
				return err
			}

			defer fh.Close()

			_, err = fh.Write([]byte(msg))

			if err != nil {
				return err
			}

			return nil
		}

	} else {
		log.Fatal("Invalid callback")
	}

	l, err := listener.NewListener(*endpoint, f)

	if err != nil {
		log.Fatal(err)
	}

	err = l.Start()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
