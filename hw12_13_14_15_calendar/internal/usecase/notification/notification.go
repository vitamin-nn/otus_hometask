package notification

import (
	"context"
	"time"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

type UseCase struct {
	er repository.EventRepo
}

func NewNotifyUseCase(repo repository.EventRepo) *UseCase {
	return &UseCase{
		er: repo,
	}
}

func (n UseCase) GetNotifications(ctx context.Context) ([]*repository.Notification, error) {
	notifications, err := n.er.GetNotifyEvents(ctx, time.Now())
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (n UseCase) SetSent(ctx context.Context, eventID int) error {
	return n.er.UpdateNotificationSent(ctx, eventID, time.Now())
}

func (n UseCase) CleanOlEvents(ctx context.Context, evDaysLive int) error {
	t := time.Now().AddDate(0, 0, -evDaysLive)

	return n.er.DeleteOldEvents(ctx, t)
}
