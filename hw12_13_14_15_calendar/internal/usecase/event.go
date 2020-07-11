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

func (e *EventUseCase) UpdateEvent(ctx context.Context, title, descr string, startAt, endAt, notifyAt time.Time, eventID, userID int) (*repository.Event, error) {
	event := &repository.Event{
		ID:          eventID,
		Title:       title,
		Description: descr,
		StartAt:     startAt,
		EndAt:       endAt,
		NotifyAt:    notifyAt,
		UserID:      userID,
	}

	return e.er.UpdateEvent(ctx, event)
}

func (e *EventUseCase) DeleteEvent(ctx context.Context, eventID int) error {
	return e.er.DeleteEvent(ctx, eventID)
}

func (e *EventUseCase) GetEventsDay(ctx context.Context, userID int, dBegin time.Time) ([]*repository.Event, error) {
	year, month, day := dBegin.Date()
	loc := dBegin.Location()
	begin := time.Date(year, month, day, 0, 0, 0, 0, loc)
	end := time.Date(year, month, day, 23, 59, 59, 0, loc)
	return e.er.GetEventsByFilter(ctx, userID, begin, end)
}

func (e *EventUseCase) GetEventsWeek(ctx context.Context, userID int, wBegin time.Time) ([]*repository.Event, error) {
	end := wBegin.AddDate(0, 0, 7)
	return e.er.GetEventsByFilter(ctx, userID, wBegin, end)
}

func (e *EventUseCase) GetEventsMonth(ctx context.Context, userID int, mBegin time.Time) ([]*repository.Event, error) {
	end := mBegin.AddDate(0, 1, 0)
	return e.er.GetEventsByFilter(ctx, userID, mBegin, end)
}
