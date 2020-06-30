package server

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/server/handler"
)

type App struct {
	httpServer *http.Server
	cHandler   *handler.Calendar
}

func NewApp(repo repository.EventsRepo) *App {
	a := new(App)
	a.cHandler = handler.NewCalendarHandler(repo)
	return a
}

func (a *App) Run(addr string, wTimeout, rTimeout time.Duration) {
	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", a.cHandler.HelloWorld)

	siteHandler := logMiddleware(siteMux)

	a.httpServer = &http.Server{
		Addr:         addr,
		WriteTimeout: wTimeout,
		ReadTimeout:  rTimeout,
		Handler:      siteHandler,
	}

	if err := a.httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
