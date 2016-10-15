package main

/*

	THIS DOESN'T REALLY DO ANYTHING. BY DESIGN. NO.

	Think of this as "starter" code for most basic HTTP pony
	servers. As in copy it and rename it as something else
	and modify the stuff that goes on inside the 'request_handler'
	function defintion.

	(20160615/thisisaaronland)

*/

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-httpony/cors"
	"github.com/whosonfirst/go-httpony/tls"
	"net/http"
	"os"
)

func main() {

	var host = flag.String("host", "localhost", "Hostname to listen on")
	var port = flag.Int("port", 8080, "Port to listen on")
	var cors_enable = flag.Bool("cors", false, "...")
	var cors_allow = flag.String("cors-allow", "*", "...")
	var tls_enable = flag.Bool("tls", false, "Serve requests over TLS") // because CA warnings in browsers...
	var tls_cert = flag.String("tls-cert", "", "Path to an existing TLS certificate. If absent a self-signed certificate will be generated.")
	var tls_key = flag.String("tls-key", "", "Path to an existing TLS key. If absent a self-signed key will be generated.")

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	request_handler := func() http.Handler {

		fn := func(rsp http.ResponseWriter, req *http.Request) {

			query := req.URL.Query()

			js, err := json.Marshal(query)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			rsp.Header().Set("Content-Type", "application/json")
			rsp.Write(js)
		}

		return http.HandlerFunc(fn)
	}

	handler := cors.EnsureCORSHandler(request_handler(), *cors_enable, *cors_allow)

	var err error

	if *tls_enable {

		var cert string
		var key string

		if *tls_cert == "" && *tls_key == "" {

			root, err := tls.EnsureTLSRoot()

			if err != nil {
				panic(err)
			}

			cert, key, err = tls.GenerateTLSCert(*host, root)

			if err != nil {
				panic(err)
			}

		} else {
			cert = *tls_cert
			key = *tls_key
		}

		fmt.Printf("start and listen for requests at https://%s\n", endpoint)
		err = http.ListenAndServeTLS(endpoint, cert, key, handler)

	} else {

		fmt.Printf("start and listen for requests at http://%s\n", endpoint)
		err = http.ListenAndServe(endpoint, handler)
	}

	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
