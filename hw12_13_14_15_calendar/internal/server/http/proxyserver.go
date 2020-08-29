package http

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	grpcServ "github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/server/grpc"
	"google.golang.org/grpc"
)

type ProxyServer struct {
	httpProxy *http.Server
}

func NewProxyServer() *ProxyServer {
	s := new(ProxyServer)

	return s
}

func (s *ProxyServer) Run(ctx context.Context, httpAddr, grpcAddr string, wTimeout, rTimeout time.Duration) {
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := grpcServ.RegisterCalendarServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		log.Fatal(err)
	}

	siteMux := http.NewServeMux()
	siteMux.Handle("/", gwmux)
	siteHandler := userIDMiddleware(logMiddleware(siteMux))

	s.httpProxy = &http.Server{
		Addr:         httpAddr,
		WriteTimeout: wTimeout,
		ReadTimeout:  rTimeout,
		Handler:      siteHandler,
	}

	if err := s.httpProxy.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *ProxyServer) Shutdown(ctx context.Context) error {
	return s.httpProxy.Shutdown(ctx)
}
