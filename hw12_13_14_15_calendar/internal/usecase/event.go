package usecase

import (
	"context"
	"time"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

type EventUseCase struct {
	er repository.EventRepo
}

func NewEventUseCase(eventRepo repository.EventRepo) *EventUseCase {
	return &EventUseCase{
		er: eventRepo,
	}
}

func (e *EventUseCase) CreateEvent(ctx context.Context, title, descr string, startAt, endAt, notifyAt time.Time, userID int) (*repository.Event, error) {
	event := &repository.Event{
		Title:       title,
		Description: descr,
		StartAt:     startAt,
		EndAt:       endAt,
		NotifyAt:    notifyAt,
		UserID:      userID,
	}

	return e.er.CreateEvent(ctx, event)
}
