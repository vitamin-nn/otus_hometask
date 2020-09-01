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
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/sender"
)

func senderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sender",
		Short: "Start sender",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadSender()
			if err != nil {
				log.Fatalf("sender config file read error: %v", err)
			}

			err = logger.Init(cfg.Log)
			if err != nil {
				log.Fatalf("initialize logger error: %v", err)
			}

			rmq := rabbit.NewConsume(cfg.Rabbit.GetDSN(), "calendar_sender", cfg.Rabbit.ExchangeName, cfg.Rabbit.ExchangeType, cfg.Rabbit.QueueName)
			err = rmq.Handle(sender.GetSenderFunc(), 1)
			if err != nil {
				log.Fatalf("consumer handle error: %v", err)
			}
			defer rmq.Close()

			interruptCh := make(chan os.Signal, 1)
			signal.Notify(interruptCh, os.Interrupt)
			log.Infof("graceful shutdown: %v", <-interruptCh)

			_, finish := context.WithTimeout(context.Background(), 5*time.Second)
			defer finish()

			log.Info("finished main program")
		},
	}
}
