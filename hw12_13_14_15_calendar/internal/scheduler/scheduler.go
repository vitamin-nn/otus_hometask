package scheduler

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/rmq"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/usecase/notification"
)

type Scheduler struct {
	rmqP rmq.Publisher
	uc   *notification.UseCase
}

func NewScheduler(rmqP rmq.Publisher, uc *notification.UseCase) *Scheduler {
	return &Scheduler{
		rmqP: rmqP,
		uc:   uc,
	}
}

func (s *Scheduler) Run(ctx context.Context, checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.process(ctx)
			if err != nil {
				log.Errorf("process schedule error: %v", err)
			}
		}
	}
}

func (s *Scheduler) process(ctx context.Context) error {
	notifications, err := s.uc.GetNotifications(ctx)
	if err != nil {
		return err
	}
	for _, n := range notifications {
		err = s.rmqP.Publish(n, "")
		if err != nil {
			log.Errorf("Publish to RMQ error: %v", err)

			continue
		}
		err = s.uc.SetSent(ctx, n.EventID)
		if err != nil {
			log.Errorf("Set as sent error: %v", err)
		}
	}

	return nil
}

func (s *Scheduler) RunCleaner(ctx context.Context, cleanInterval time.Duration, evDaysLive int) {
	ticker := time.NewTicker(cleanInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.uc.CleanOlEvents(ctx, evDaysLive)
			if err != nil {
				log.Errorf("delete old events error: %v", err)
			}
		}
	}
}
