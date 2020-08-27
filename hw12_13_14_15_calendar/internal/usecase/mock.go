package usecase

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

type EventUseCaseMock struct {
	mock.Mock
}

func (m *EventUseCaseMock) CreateEvent(ctx context.Context, title, descr string, startAt, endAt, notifyAt time.Time, userID int) (*repository.Event, error) {
	args := m.Called(title, descr, startAt, endAt, notifyAt, userID)

	return args.Get(0).(*repository.Event), args.Error(1)
}

func (m *EventUseCaseMock) UpdateEvent(ctx context.Context, title, descr string, startAt, endAt, notifyAt time.Time, eventID, userID int) (*repository.Event, error) {
	args := m.Called(title, descr, startAt, endAt, notifyAt, eventID, userID)

	return args.Get(0).(*repository.Event), args.Error(1)
}

func (m *EventUseCaseMock) DeleteEvent(ctx context.Context, eventID int) error {
	args := m.Called(eventID)

	return args.Error(0)
}

func (m *EventUseCaseMock) GetEventsDay(ctx context.Context, userID int, dBegin time.Time) ([]*repository.Event, error) {
	args := m.Called(userID, dBegin)

	return args.Get(0).([]*repository.Event), args.Error(1)
}

func (m *EventUseCaseMock) GetEventsWeek(ctx context.Context, userID int, wBegin time.Time) ([]*repository.Event, error) {
	args := m.Called(userID, wBegin)

	return args.Get(0).([]*repository.Event), args.Error(1)
}

func (m *EventUseCaseMock) GetEventsMonth(ctx context.Context, userID int, mBegin time.Time) ([]*repository.Event, error) {
	args := m.Called(userID, mBegin)

	return args.Get(0).([]*repository.Event), args.Error(1)
}
