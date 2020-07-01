package inmemory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

func TestEvent(t *testing.T) {
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
		{
			ID:          1,
			Title:       "Test event 1",
			Description: "Test event 1 description",
			StartAt:     t1Start,
			EndAt:       t1End,
			UserID:      userID,
		},
		{
			ID:          2,
			Title:       "Test event 2",
			Description: "Test event 2 description",
			StartAt:     t2Start,
			EndAt:       t2End,
			NotifyAt:    t2Notify,
			UserID:      userID,
		},
		{
			ID:          3,
			Title:       "Test event 3",
			Description: "Test event 3 description",
			StartAt:     t3Start,
			EndAt:       t3End,
			UserID:      userID,
		},
	}

	t4Start, err := time.Parse(time.RFC3339, "2020-01-02T15:30:00+03:00")
	require.Nil(t, err)
	t4End, err := time.Parse(time.RFC3339, "2020-01-02T16:30:00+03:00")
	require.Nil(t, err)
	overlapEvent := &repository.Event{
		ID:          4,
		Title:       "Test event 4",
		Description: "Test event 4 description",
		StartAt:     t4Start,
		EndAt:       t4End,
		UserID:      userID,
	}
	ctx := context.Background()

	t.Run("create+get events", func(t *testing.T) {
		repo := NewEventRepo()
		for _, event := range events {
			err := repo.CreateEvent(ctx, event)
			require.Nil(t, err)
		}
		evt, err := repo.GetEventByID(ctx, 1)
		require.Nil(t, err)
		require.Equal(t, evt.ID, 1)
		require.Equal(t, evt.Title, "Test event 1")

		err = repo.CreateEvent(ctx, overlapEvent)
		require.Equal(t, repository.ErrDateBusy, err)
	})

	t.Run("update events", func(t *testing.T) {
		repo := NewEventRepo()
		eventID := 1
		for _, event := range events {
			_ = repo.CreateEvent(ctx, event)
		}
		// Get тестируется в другом методе
		evt, _ := repo.GetEventByID(ctx, eventID)
		title := "Updated test event 1"
		evt.Title = title
		startAt, _ := time.Parse(time.RFC3339, "2020-01-02T14:00:00+03:00")
		evt.StartAt = startAt
		err := repo.UpdateEvent(ctx, eventID, evt)
		require.Nil(t, err)
		evt, _ = repo.GetEventByID(ctx, eventID)
		require.Equal(t, evt.Title, title)
		require.Equal(t, evt.StartAt, startAt)

		evt = overlapEvent
		evt.ID = eventID
		err = repo.UpdateEvent(ctx, eventID, evt)
		require.Equal(t, repository.ErrDateBusy, err)
	})

	t.Run("delete events", func(t *testing.T) {
		repo := NewEventRepo()
		eventID := 1
		unknownEventID := 10
		for _, event := range events {
			_ = repo.CreateEvent(ctx, event)
		}
		err := repo.DeleteEvent(ctx, unknownEventID)
		require.Equal(t, repository.ErrEventNotFound, err)
		err = repo.DeleteEvent(ctx, eventID)
		require.Nil(t, err)
		_, err = repo.GetEventByID(ctx, eventID)
		require.Equal(t, repository.ErrEventNotFound, err)
	})

	t.Run("get day events", func(t *testing.T) {
		repo := NewEventRepo()
		for _, event := range events {
			_ = repo.CreateEvent(ctx, event)
		}
		dBegin, err := time.Parse(time.RFC3339, "2020-01-02T00:00:00+03:00")
		require.Nil(t, err)

		events, err := repo.GetEventsDay(ctx, userID, dBegin)
		require.Nil(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, events[0].ID, 1)
	})

	t.Run("get week events", func(t *testing.T) {
		repo := NewEventRepo()
		for _, event := range events {
			_ = repo.CreateEvent(ctx, event)
		}
		wBegin, err := time.Parse(time.RFC3339, "2019-12-30T00:00:00+03:00")
		require.Nil(t, err)

		events, err := repo.GetEventsWeek(ctx, userID, wBegin)
		require.Nil(t, err)
		require.Equal(t, 2, len(events))
	})

	t.Run("get month events", func(t *testing.T) {
		repo := NewEventRepo()
		for _, event := range events {
			_ = repo.CreateEvent(ctx, event)
		}
		mBegin, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00+03:00")
		require.Nil(t, err)

		events, err := repo.GetEventsMonth(ctx, userID, mBegin)
		require.Nil(t, err)
		require.Equal(t, 3, len(events))
	})
}
