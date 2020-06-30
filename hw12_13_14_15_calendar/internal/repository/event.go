package repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("time is busy")
)

type Event struct {
	ID          int
	Title       string
	Description string
	StartAt     time.Time
	EndAt       time.Time
	NotifyAt    time.Time
	UserID      int
}

type Notification struct {
	EventID      int
	EventTitle   string
	StartAt      time.Time
	NotifyUserID int
}

type EventsRepo interface {
	CreateEvent(ctx context.Context, event *Event) error
	UpdateEvent(ctx context.Context, eventID int, event *Event) error
	DeleteEvent(ctx context.Context, eventID int) error
	GetEventsDay(ctx context.Context, userID int, dBegin time.Time) ([]*Event, error)
	GetEventsWeek(ctx context.Context, userID int, wBegin time.Time) ([]*Event, error)
	GetEventsMonth(ctx context.Context, userID int, mBegin time.Time) ([]*Event, error)
	GetEventByID(ctx context.Context, eventID int) (*Event, error)
}
