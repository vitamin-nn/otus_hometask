package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/config"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository/inmemory"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository/psql"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "calendar",
		Short: "Calendar service",
	}

	rootCmd.AddCommand(serverCmd())
	rootCmd.AddCommand(schedulerCmd())
	rootCmd.AddCommand(senderCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute cmd: %v", err)
	}
}

func getRepo(ctx context.Context, repoType string, cpg config.PSQL) repository.EventRepo {
	var repo repository.EventRepo
	switch repoType {
	case inmemory.Type:
		repo = inmemory.NewEventRepo()
	case psql.Type:
		repo = psql.NewEventRepo(cpg.GetDSN())
		err := repo.Connect(ctx)
		if err != nil {
			log.Fatalf("psql connect error: %v", err)
		}
	}

	return repo
}
