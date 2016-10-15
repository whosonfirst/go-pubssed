package cors

import (
	"net/http"
)

func EnsureCORSHandler(next http.Handler, enable bool, allow string) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		if enable {
			rsp.Header().Set("Access-Control-Allow-Origin", allow)
		}

		next.ServeHTTP(rsp, req)
	}

	return http.HandlerFunc(fn)
}
