package server

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s %v %s",
			r.Method,
			r.Proto,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
			r.UserAgent(),
		)
	})
}
