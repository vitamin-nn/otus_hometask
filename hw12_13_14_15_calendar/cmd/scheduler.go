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
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/rmq/rabbit"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/usecase/notification"
)

func schedulerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scheduler",
		Short: "Start scheduler",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadScheduler()
			if err != nil {
				log.Fatalf("scheduler config file read error: %v", err)
			}

			err = logger.Init(cfg.Log)
			if err != nil {
				log.Fatalf("initialize logger error: %v", err)
			}

			rmq := rabbit.NewPublish(cfg.Rabbit.GetDSN(), cfg.Rabbit.ExchangeName, cfg.Rabbit.ExchangeType, cfg.Rabbit.QueueName)
			err = rmq.Connect()
			if err != nil {
				log.Fatalf("rabbit connect error: %v", err)
			}
			defer rmq.Close()

			ctx, cancel := context.WithCancel(context.Background())
			repo := getRepo(ctx, cfg.App.RepoType, cfg.PSQL)
			defer repo.Close()
			uc := notification.NewNotifyUseCase(repo)

			sc := scheduler.NewScheduler(rmq, uc)

			go func() {
				sc.Run(ctx, cfg.Scheduler.CheckInterval)
			}()

			go func() {
				sc.RunCleaner(ctx, cfg.Scheduler.CleanInterval, cfg.Scheduler.EventLiveDays)
			}()

			interruptCh := make(chan os.Signal, 1)
			signal.Notify(interruptCh, os.Interrupt)
			log.Infof("graceful shutdown: %v", <-interruptCh)

			_, finish := context.WithTimeout(context.Background(), 5*time.Second)
			defer finish()
			cancel()

			log.Info("finished main program")
		},
	}
}
