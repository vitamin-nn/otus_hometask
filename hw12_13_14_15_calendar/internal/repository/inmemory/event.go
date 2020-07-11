package inmemory

import (
	"context"
	"errors"
	"sync"
	"time"

	outErr "github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/error"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

const Type = "inmemory"

var ErrDontUse = errors.New("don't use this method for in_memory implementation")

var _ repository.EventRepo = (*InMemory)(nil)

type InMemory struct {
	eventsCounter int
	events        map[int]*repository.Event
	mutex         *sync.RWMutex
}

func NewEventRepo() *InMemory {
	return &InMemory{
		events: make(map[int]*repository.Event),
		mutex:  new(sync.RWMutex),
	}
}

func (e *InMemory) CreateEvent(ctx context.Context, event *repository.Event) (*repository.Event, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	isBusy, err := e.isBusyTime(ctx, event.UserID, event.StartAt, event.EndAt)
	if err != nil {
		return nil, err
	}
	if isBusy {
		return nil, outErr.ErrDateBusy
	}
	e.eventsCounter++
	event.ID = e.eventsCounter
	e.events[event.ID] = event
	return event, nil
}

func (e *InMemory) UpdateEvent(ctx context.Context, event *repository.Event) (*repository.Event, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	isBusy, err := e.isBusyTime(ctx, event.UserID, event.StartAt, event.EndAt)
	if err != nil {
		return nil, err
	}
	if isBusy {
		return nil, outErr.ErrDateBusy
	}
	e.events[event.ID] = event
	return event, nil
}

func (e *InMemory) DeleteEvent(_ context.Context, eventID int) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	_, ok := e.events[eventID]
	if ok {
		delete(e.events, eventID)
		return nil
	}

	return outErr.ErrEventNotFound
}

func (e *InMemory) GetEventByID(_ context.Context, eventID int) (*repository.Event, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	evt, ok := e.events[eventID]
	if ok {
		return evt, nil
	}

	return nil, outErr.ErrEventNotFound
}

func (e *InMemory) GetEventsByFilter(ctx context.Context, userID int, begin time.Time, end time.Time) ([]*repository.Event, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.getEventsByFilterInternal(ctx, userID, begin, end)
}

func (e *InMemory) getEventsByFilterInternal(_ context.Context, userID int, begin time.Time, end time.Time) ([]*repository.Event, error) {
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

func (e *InMemory) isBusyTime(ctx context.Context, userID int, begin time.Time, end time.Time) (bool, error) {
	events, err := e.getEventsByFilterInternal(ctx, userID, begin, end)
	if err != nil {
		return false, err
	}
	if len(events) == 0 {
		return false, nil
	}
	return true, nil
}

func (e *InMemory) Connect(_ context.Context) error {
	// возможно, есть варианты лучше как поступать в подобных ситуация
	return ErrDontUse
}

func (e *InMemory) Close() error {
	return ErrDontUse
}
