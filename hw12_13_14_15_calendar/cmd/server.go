package cmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/config"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/logger"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/server/http"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/usecase/calendar"
)

func serverCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start calendar grpc and http servers",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadServer()
			if err != nil {
				log.Fatalf("server config file read error: %v", err)
			}

			err = logger.Init(cfg.Log)
			if err != nil {
				log.Fatalf("initialize logger error: %v", err)
			}

			repo := getRepo(context.Background(), cfg.App.RepoType, cfg.PSQL)
			defer repo.Close()

			eUseCase := calendar.NewEventUseCase(repo)
			grpcServer := grpc.NewCalendarServer(eUseCase)
			proxyServer := http.NewProxyServer()

			log.WithFields(cfg.Fields()).Info("Starting calendar service")

			ctx := context.Background()
			go func() {
				log.Info("Starting GRPC server")
				if err := grpcServer.Run(cfg.GrpcServer.Addr); err != nil {
					log.Fatal(err)
				}
			}()
			go func() {
				log.Info("Starting HTTP Proxy server")
				proxyServer.Run(ctx, cfg.Server.Addr, cfg.GrpcServer.Addr, cfg.Server.WriteTimeout, cfg.Server.ReadTimeout)
			}()

			interruptCh := make(chan os.Signal, 1)
			signal.Notify(interruptCh, os.Interrupt)
			log.Infof("graceful shutdown: %v", <-interruptCh)

			ctx, finish := context.WithTimeout(context.Background(), 5*time.Second)
			defer finish()
			grpcServer.Down()
			err = proxyServer.Shutdown(ctx)
			if err != nil {
				log.Errorf("error while shutdown: %v", err)
			}

			log.Info("finished main program")
		},
	}
}
