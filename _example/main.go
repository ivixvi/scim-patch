package main

import (
	logger "log"
	"net/http"
)

func main() {
	server := newTestServer()
	loggingMiddleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Println(r.Method, r.URL.Path)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
	logger.Fatal(http.ListenAndServe(":7643", loggingMiddleware(server)))
}
