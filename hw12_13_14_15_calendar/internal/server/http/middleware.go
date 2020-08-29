package http

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/server"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)

		defer log.Printf(
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

func userIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(server.UserIDHeaderKey)
		log.Infof("User_id: %s", userID)
		ctx := context.WithValue(r.Context(), server.UserIDGrpcKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
