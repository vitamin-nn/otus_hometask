package repository

import (
	"context"
	"time"
)

type Event struct {
	ID               int
	Title            string
	Description      string
	StartAt          time.Time
	EndAt            time.Time
	NotifyAt         time.Time
	NotificationSent time.Time
	UserID           int
}

type Notification struct {
	EventID      int
	EventTitle   string
	StartAt      time.Time
	NotifyUserID int
}

type EventRepo interface {
	Connect(ctx context.Context) error
	Close() error
	CreateEvent(ctx context.Context, event *Event) (*Event, error)
	UpdateEvent(ctx context.Context, event *Event) (*Event, error)
	DeleteEvent(ctx context.Context, eventID int) error
	GetEventsByFilter(ctx context.Context, userID int, begin time.Time, end time.Time) ([]*Event, error)
	GetEventByID(ctx context.Context, eventID int) (*Event, error)
	DeleteOldEvents(ctx context.Context, t time.Time) error
	GetNotifyEvents(ctx context.Context, t time.Time) ([]*Notification, error)
	UpdateNotificationSent(ctx context.Context, eventID int, t time.Time) error
}
