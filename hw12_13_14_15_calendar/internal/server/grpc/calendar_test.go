package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	outErr "github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/error"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/server"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/usecase/calendar"
	"google.golang.org/grpc/metadata"
)

func TestGrpcModifyEvent(t *testing.T) {
	uc := new(calendar.EventUseCaseMock)
	s := NewCalendarServer(uc)

	tStart, err := time.Parse(time.RFC3339, "2020-01-02T15:00:00Z")
	require.Nil(t, err)
	tStartProto, err := ptypes.TimestampProto(tStart)
	require.Nil(t, err)

	tEnd, err := time.Parse(time.RFC3339, "2020-01-02T16:00:00Z")
	require.Nil(t, err)
	tEndProto, err := ptypes.TimestampProto(tEnd)
	require.Nil(t, err)

	tNotify, err := time.Parse(time.RFC3339, "2020-01-02T14:00:00Z")
	require.Nil(t, err)
	tNotifyProto, err := ptypes.TimestampProto(tNotify)
	require.Nil(t, err)

	id := 1
	title := "test"
	descr := "test"
	userID := 1
	eProto := &Event{
		Id:          int32(id),
		Title:       title,
		Description: descr,
		StartAt:     tStartProto,
		EndAt:       tEndProto,
		NotifyAt:    tNotifyProto,
		UserId:      int32(userID),
	}
	eRepo := &repository.Event{
		ID:          id,
		Title:       title,
		Description: descr,
		StartAt:     tStart,
		EndAt:       tEnd,
		NotifyAt:    tNotify,
		UserID:      userID,
	}

	eReq := &ModifyEventRequest{
		Title:       title,
		Description: descr,
		StartAt:     tStartProto,
		EndAt:       tEndProto,
		NotifyAt:    tNotifyProto,
	}

	t.Run("create event", func(t *testing.T) {
		uc.On("CreateEvent", title, descr, tStart, tEnd, tNotify, userID).Return(eRepo, nil)

		req := &CreateEventRequest{
			Event: eReq,
		}

		resp, err := s.CreateEvent(getGrpcCtx(), req)
		require.Nil(t, err)

		res, ok := resp.Result.(*ModifyEventResponse_Event)
		require.True(t, ok)
		require.Equal(t, eProto, res.Event)
	})

	t.Run("update event", func(t *testing.T) {
		title = "test_upd"
		eProto.Title = title
		eRepo.Title = title
		eReq.Title = title
		uc.On("UpdateEvent", title, descr, tStart, tEnd, tNotify, id, userID).Return(eRepo, nil)

		req := &UpdateEventRequest{
			Event:   eReq,
			EventId: int32(id),
		}

		resp, err := s.UpdateEvent(getGrpcCtx(), req)
		require.Nil(t, err)

		res, ok := resp.Result.(*ModifyEventResponse_Event)
		require.True(t, ok)
		require.Equal(t, eProto, res.Event)
	})

	t.Run("error modify event", func(t *testing.T) {
		ucErr := new(calendar.EventUseCaseMock)
		sErr := NewCalendarServer(ucErr)

		// в коде реализации библиотеки Mock в методе On происходит append для каждого нового метода,
		// поэтому нельзя переопределить уже существующий UpdateEvent и пришлось создать отдельный объект ucErr
		ucErr.On("UpdateEvent", title, descr, tStart, tEnd, tNotify, id, userID).Return(eRepo, errors.New("Unknown error"))

		req := &UpdateEventRequest{
			Event:   eReq,
			EventId: int32(id),
		}

		resp, err := sErr.UpdateEvent(getGrpcCtx(), req)
		require.Nil(t, err)

		res, ok := resp.Result.(*ModifyEventResponse_Error)
		require.True(t, ok)
		require.Equal(t, outErr.ErrInternal.Error(), res.Error)
	})
}

func TestGrpcDeleteEvent(t *testing.T) {
	t.Run("delete event", func(t *testing.T) {
		uc := new(calendar.EventUseCaseMock)
		s := NewCalendarServer(uc)

		id := 1
		uc.On("DeleteEvent", id).Return(nil)

		req := &DeleteEventRequest{
			EventId: int32(id),
		}
		resp, err := s.DeleteEvent(getGrpcCtx(), req)
		require.Nil(t, err)

		res, ok := resp.Result.(*DeleteResponse_Success)
		require.True(t, ok)
		require.True(t, res.Success)
	})
	t.Run("delete event with error", func(t *testing.T) {
		uc := new(calendar.EventUseCaseMock)
		s := NewCalendarServer(uc)

		retErr := outErr.ErrEventNotFound
		id := 1
		uc.On("DeleteEvent", id).Return(retErr)

		req := &DeleteEventRequest{
			EventId: int32(id),
		}
		resp, err := s.DeleteEvent(getGrpcCtx(), req)
		require.Nil(t, err)

		res, ok := resp.Result.(*DeleteResponse_Error)
		require.True(t, ok)
		require.Equal(t, retErr.Error(), res.Error)
	})
}

func TestGrpcFilterEvent(t *testing.T) {
	uc := new(calendar.EventUseCaseMock)
	s := NewCalendarServer(uc)

	// здесь нет необходиомсти задавать корректные даты, т.к. используем mock
	// а правильность бизнес-логики проверяется в других тестах
	title := "test"
	descr := "test"
	userID := 1
	timeRepo, err := time.Parse(time.RFC3339, "2020-01-02T15:00:00Z")
	require.Nil(t, err)
	timeProto, err := ptypes.TimestampProto(timeRepo)
	require.Nil(t, err)

	req := &GetEventsRequest{
		BeginAt: timeProto,
	}

	var protoEvents []*Event
	var repoEvents []*repository.Event

	for i := 1; i <= 3; i++ {
		eProto := &Event{
			Id:          int32(i),
			Title:       title,
			Description: descr,
			StartAt:     timeProto,
			EndAt:       timeProto,
			NotifyAt:    timeProto,
			UserId:      int32(userID),
		}
		protoEvents = append(protoEvents, eProto)
		eRepo := &repository.Event{
			ID:          i,
			Title:       title,
			Description: descr,
			StartAt:     timeRepo,
			EndAt:       timeRepo,
			NotifyAt:    timeRepo,
			UserID:      userID,
		}
		repoEvents = append(repoEvents, eRepo)
	}

	checkFunc := func(resp *GetEventsResponse, err error) {
		require.Nil(t, err)

		res, ok := resp.Result.(*GetEventsResponse_Events)
		require.True(t, ok)

		require.Equal(t, len(protoEvents), len(res.Events.Events))
		for k, e := range res.Events.Events {
			require.Equal(t, protoEvents[k], e)
		}
	}

	t.Run("day events", func(t *testing.T) {
		uc.On("GetEventsDay", userID, timeRepo).Return(repoEvents, nil)
		resp, err := s.GetEventsDay(getGrpcCtx(), req)
		checkFunc(resp, err)
	})

	t.Run("week events", func(t *testing.T) {
		uc.On("GetEventsWeek", userID, timeRepo).Return(repoEvents, nil)
		resp, err := s.GetEventsWeek(getGrpcCtx(), req)
		checkFunc(resp, err)
	})

	t.Run("month events", func(t *testing.T) {
		uc.On("GetEventsMonth", userID, timeRepo).Return(repoEvents, nil)
		resp, err := s.GetEventsMonth(getGrpcCtx(), req)
		checkFunc(resp, err)
	})
}

func getGrpcCtx() context.Context {
	m := metadata.New(map[string]string{server.UserIDHeaderKey: "1"})
	return metadata.NewIncomingContext(context.Background(), m)
}
