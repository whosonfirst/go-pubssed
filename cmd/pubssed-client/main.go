package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/whosonfirst/go-pubssed/listener"
)

func main() {

	var endpoint = flag.String("endpoint", "", "The pubssed endpoint you are connecting to.")
	var callback = flag.String("callback", "debug", "The callback to invoke when a SSE event is received.")
	var append_root = flag.String("append-root", ".", "The destination to write log files if the 'append' callback is invoked.")
	var retry = flag.Bool("retry-on-eof", false, "Try to reconnect to the SSE endpoint if an EOF error is triggered. This is sometimes necessary if an SSE endpoint is configured with a too-short HTTP timeout (for example if running behind an AWS load balancer).")

	var verbose = flag.Bool("verbose", false, "Enable verbose (debug) logging")

	flag.Parse()

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	if *endpoint == "" {
		log.Fatal("Missing pubssed endpoint")
	}

	var f listener.ListenerFunc

	if *callback == "debug" {

		f = func(msg string) error {

			log.Println(msg)
			return nil
		}

	} else if *callback == "append" {

		root := *append_root

		if root == "." {

			cwd, err := os.Getwd()

			if err != nil {
				log.Fatal(err)
			}

			root = cwd
		} else {

			_, err := os.Stat(root)

			if os.IsNotExist(err) {
				log.Fatal(root)
			}
		}

		f = func(msg string) error {

			now := time.Now()

			// ts := now.UnixNano() / int64(time.Millisecond)

			year := fmt.Sprintf("%04d", now.Year())
			month := fmt.Sprintf("%02d", now.Month())
			day := fmt.Sprintf("%02d", now.Day())

			dirname := filepath.Join(root, year, month, day)
			fname := fmt.Sprintf("%s%s%s%02d.txt", year, month, day, now.Hour())

			path := filepath.Join(dirname, fname)

			_, err := os.Stat(dirname)

			if os.IsNotExist(err) {

				err = os.MkdirAll(dirname, 0755)

				if err != nil {
					return err
				}
			}

			fh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

			if err != nil {
				return err
			}

			defer fh.Close()

			msg = fmt.Sprintf("%s\n", msg)

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

	for {
		err = l.Start()

		if err != nil {

			if *retry && err.Error() == "EOF" {
				log.Println("EOF error triggered, reconnecting")
				continue
			}

			log.Fatal(err)
		}
	}

	os.Exit(0)
}
