package calendar

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository/inmemory"
)

func TestUsecaseEvent(t *testing.T) {
	t1Start, err := time.Parse(time.RFC3339, "2020-01-02T15:00:00+03:00")
	require.Nil(t, err)
	t1End, err := time.Parse(time.RFC3339, "2020-01-02T16:00:00+03:00")
	require.Nil(t, err)
	t2Start, err := time.Parse(time.RFC3339, "2020-01-04T15:00:00+03:00")
	require.Nil(t, err)
	t2End, err := time.Parse(time.RFC3339, "2020-01-04T16:00:00+03:00")
	require.Nil(t, err)
	t2Notify, err := time.Parse(time.RFC3339, "2020-01-04T12:00:00+03:00")
	require.Nil(t, err)
	t3Start, err := time.Parse(time.RFC3339, "2020-01-30T17:00:00+03:00")
	require.Nil(t, err)
	t3End, err := time.Parse(time.RFC3339, "2020-01-30T18:00:00+03:00")
	require.Nil(t, err)
	userID := 1
	events := []*repository.Event{
		{ // ID:          1,
			Title:       "Test event 1",
			Description: "Test event 1 description",
			StartAt:     t1Start,
			EndAt:       t1End,
			UserID:      userID,
		},
		{ // ID:          2,
			Title:       "Test event 2",
			Description: "Test event 2 description",
			StartAt:     t2Start,
			EndAt:       t2End,
			NotifyAt:    t2Notify,
			UserID:      userID,
		},
		{ // ID:          3,
			Title:       "Test event 3",
			Description: "Test event 3 description",
			StartAt:     t3Start,
			EndAt:       t3End,
			UserID:      userID,
		},
	}

	makeUsecase := func() *EventUseCase {
		inMemoryRepo := inmemory.NewEventRepo()
		usecase := NewEventUseCase(inMemoryRepo)
		ctx := context.Background()
		for _, event := range events {
			_, err := usecase.CreateEvent(ctx, event.Title, event.Description, event.StartAt, event.EndAt, event.NotifyAt, event.UserID)
			require.Nil(t, err)
		}
		return usecase
	}

	ctx := context.Background()

	t.Run("get day events", func(t *testing.T) {
		uc := makeUsecase()
		dBegin, err := time.Parse(time.RFC3339, "2020-01-02T00:00:00+03:00")
		require.Nil(t, err)

		events, err := uc.GetEventsDay(ctx, userID, dBegin)
		require.Nil(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, events[0].ID, 1)
	})

	t.Run("get week events", func(t *testing.T) {
		uc := makeUsecase()
		wBegin, err := time.Parse(time.RFC3339, "2019-12-30T00:00:00+03:00")
		require.Nil(t, err)

		events, err := uc.GetEventsWeek(ctx, userID, wBegin)
		require.Nil(t, err)
		require.Equal(t, 2, len(events))
	})

	t.Run("get month events", func(t *testing.T) {
		uc := makeUsecase()
		mBegin, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00+03:00")
		require.Nil(t, err)

		events, err := uc.GetEventsMonth(ctx, userID, mBegin)
		require.Nil(t, err)
		require.Equal(t, 3, len(events))
	})
}
