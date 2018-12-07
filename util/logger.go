package util

import (
	"log"
	"net/http"
	"time"
)

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, req)
		finish := time.Now()

		log.Println(req.URL, req.Method, finish.Sub(start))
	})
}
