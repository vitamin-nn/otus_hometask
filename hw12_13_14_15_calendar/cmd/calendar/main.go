package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"os/signal"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/config"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/logger"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository/inmemory"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository/psql"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/server"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/usecase"
)

func main() {
	var cfgFileName string
	flag.StringVar(&cfgFileName, "config", "./configs/local.toml", "config filepath")
	cfg, err := config.Load(cfgFileName)
	if err != nil {
		log.Fatalf("config file read error: %v", err)
	}

	err = logger.Init(cfg.Log)
	if err != nil {
		log.Fatalf("initialize logger error: %v", err)
	}

	wTimeout, err := time.ParseDuration(cfg.Server.WriteTimeout)
	if err != nil {
		log.Fatalf("write timeout parse error: %v", err)
	}
	rTimeout, err := time.ParseDuration(cfg.Server.ReadTimeout)
	if err != nil {
		log.Fatalf("read timeout parse error: %v", err)
	}

	ctx := context.Background()
	var repo repository.EventRepo
	switch cfg.App.RepoType {
	case inmemory.Type:
		repo = inmemory.NewEventRepo()
	case psql.Type:
		db, err := psqlConnect(ctx, cfg.PSQL.DSN)
		if err != nil {
			log.Fatalf("psql connect error: %v", err)
		}
		defer db.Close()
		repo = psql.NewEventRepo(db)
	}

	eUseCase := usecase.NewEventUseCase(repo)
	app := server.NewApp(eUseCase)
	go app.Run(cfg.Server.Addr, wTimeout, rTimeout)

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	log.Infof("graceful shutdown: %v", <-interruptCh)

	ctx, finish := context.WithTimeout(context.Background(), 5*time.Second)
	defer finish()
	err = app.Shutdown(ctx)
	if err != nil {
		log.Errorf("error while shutdown: %v", err)
	}

	log.Info("finished main program")
}

func psqlConnect(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.Stats()
	return db, db.PingContext(ctx)
}
