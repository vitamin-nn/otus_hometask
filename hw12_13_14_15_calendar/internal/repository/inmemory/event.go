package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

const Type = "inmemory"

var _ repository.EventsRepo = (*InMemory)(nil)

type InMemory struct {
	events map[int]*repository.Event
	mutex  *sync.Mutex
}

func NewEventRepo() *InMemory {
	return &InMemory{
		events: make(map[int]*repository.Event),
		mutex:  new(sync.Mutex),
	}
}

func (e *InMemory) CreateEvent(ctx context.Context, event *repository.Event) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.isBusyTime(ctx, event.UserID, event.StartAt, event.EndAt) {
		return repository.ErrDateBusy
	}
	e.events[event.ID] = event
	return nil
}

func (e *InMemory) UpdateEvent(ctx context.Context, eventID int, event *repository.Event) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.isBusyTime(ctx, event.UserID, event.StartAt, event.EndAt) {
		return repository.ErrDateBusy
	}
	e.events[eventID] = event
	return nil
}

func (e *InMemory) DeleteEvent(_ context.Context, eventID int) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	_, ok := e.events[eventID]
	if ok {
		delete(e.events, eventID)
		return nil
	}

	return repository.ErrEventNotFound
}

func (e *InMemory) GetEventByID(_ context.Context, eventID int) (*repository.Event, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	evt, ok := e.events[eventID]
	if ok {
		return evt, nil
	}

	return nil, repository.ErrEventNotFound
}

func (e *InMemory) GetEventsDay(ctx context.Context, userID int, dBegin time.Time) ([]*repository.Event, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	end := dBegin.AddDate(0, 0, 1)
	return e.getEventsByFilter(ctx, userID, dBegin, end)
}

func (e *InMemory) GetEventsWeek(ctx context.Context, userID int, wBegin time.Time) ([]*repository.Event, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	end := wBegin.AddDate(0, 0, 7)
	return e.getEventsByFilter(ctx, userID, wBegin, end)
}

func (e *InMemory) GetEventsMonth(ctx context.Context, userID int, mBegin time.Time) ([]*repository.Event, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	end := mBegin.AddDate(0, 1, 0)
	return e.getEventsByFilter(ctx, userID, mBegin, end)
}

func (e *InMemory) getEventsByFilter(_ context.Context, userID int, begin time.Time, end time.Time) ([]*repository.Event, error) {
	var result []*repository.Event
	for _, ev := range e.events {
		if ev.UserID != userID {
			continue
		}
		if (begin.After(ev.StartAt) && begin.Before(ev.EndAt)) ||
			(end.After(ev.EndAt) && end.Before(ev.EndAt)) ||
			(begin.Before(ev.StartAt) && end.After(ev.EndAt)) {
			result = append(result, ev)
		}
	}
	return result, nil
}

func (e *InMemory) isBusyTime(_ context.Context, userID int, begin time.Time, end time.Time) bool {
	for _, ev := range e.events {
		if ev.UserID != userID {
			continue
		}
		if (begin.After(ev.StartAt) && begin.Before(ev.EndAt)) ||
			(end.After(ev.StartAt) && end.Before(ev.EndAt)) ||
			(begin.Before(ev.StartAt) && end.After(ev.EndAt)) {
			return true
		}
	}
	return false
}
