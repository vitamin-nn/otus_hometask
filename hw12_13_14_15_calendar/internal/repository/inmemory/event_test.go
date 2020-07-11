package inmemory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	outErr "github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/error"
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

	makeRepo := func() *InMemory {
		repo := NewEventRepo()
		ctx := context.Background()
		for _, event := range events {
			_, err := repo.CreateEvent(ctx, event)
			require.Nil(t, err)
		}
		return repo
	}

	t4Start, err := time.Parse(time.RFC3339, "2020-01-02T15:30:00+03:00")
	require.Nil(t, err)
	t4End, err := time.Parse(time.RFC3339, "2020-01-02T16:30:00+03:00")
	require.Nil(t, err)
	overlapEvent := &repository.Event{
		Title:       "Test event 4",
		Description: "Test event 4 description",
		StartAt:     t4Start,
		EndAt:       t4End,
		UserID:      userID,
	}
	ctx := context.Background()

	t.Run("create+get events", func(t *testing.T) {
		repo := makeRepo()
		eventID := 1
		evt, err := repo.GetEventByID(ctx, eventID)
		require.Nil(t, err)
		require.Equal(t, evt.ID, eventID)
		require.Equal(t, evt.Title, "Test event 1")

		_, err = repo.CreateEvent(ctx, overlapEvent)
		require.Equal(t, outErr.ErrDateBusy, err)
	})

	t.Run("update events", func(t *testing.T) {
		repo := makeRepo()
		eventID := 1
		// Get тестируется в другом методе
		evt, _ := repo.GetEventByID(ctx, eventID)
		title := "Updated test event 1"
		evt.Title = title
		startAt, _ := time.Parse(time.RFC3339, "2020-01-02T14:00:00+03:00")
		evt.StartAt = startAt
		evt, err := repo.UpdateEvent(ctx, evt)
		require.Nil(t, err)
		require.Equal(t, evt.Title, title)
		require.Equal(t, evt.StartAt, startAt)

		evt = overlapEvent
		evt.ID = eventID
		_, err = repo.UpdateEvent(ctx, evt)
		require.Equal(t, outErr.ErrDateBusy, err)
	})

	t.Run("delete events", func(t *testing.T) {
		repo := makeRepo()
		eventID := 1
		unknownEventID := 10
		err := repo.DeleteEvent(ctx, unknownEventID)
		require.Equal(t, outErr.ErrEventNotFound, err)
		err = repo.DeleteEvent(ctx, eventID)
		require.Nil(t, err)
		_, err = repo.GetEventByID(ctx, eventID)
		require.Equal(t, outErr.ErrEventNotFound, err)
	})

	t.Run("get events by filter", func(t *testing.T) {
		repo := makeRepo()
		begin, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00+03:00")
		require.Nil(t, err)
		end, err := time.Parse(time.RFC3339, "2020-01-15T00:00:00+03:00")
		require.Nil(t, err)

		events, err := repo.GetEventsByFilter(ctx, userID, begin, end)
		require.Nil(t, err)
		require.Equal(t, 2, len(events))
	})
}
