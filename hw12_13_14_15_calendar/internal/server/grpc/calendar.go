package grpc

import (
	"context"

	log "github.com/sirupsen/logrus"
	outErr "github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/error"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

func (s *CalendarServer) CreateEvent(ctx context.Context, req *CreateEventRequest) (*ModifyEventResponse, error) {
	return s.modifyEvent(ctx, req.GetEvent(), 0)
}

func (s *CalendarServer) UpdateEvent(ctx context.Context, req *UpdateEventRequest) (*ModifyEventResponse, error) {
	return s.modifyEvent(ctx, req.GetEvent(), int(req.GetEventId()))
}

func (s *CalendarServer) modifyEvent(ctx context.Context, eventIn *ModifyEventRequest, eventID int) (*ModifyEventResponse, error) {
	userID, err := getUser(ctx)
	if err != nil {
		return nil, err
	}

	if eventIn == nil {
		return nil, outErr.ErrEmptyEvent
	}

	startAt, err := getTimeByTimestamp(eventIn.GetStartAt())
	if err != nil {
		return nil, err
	}

	endAt, err := getTimeByTimestamp(eventIn.GetEndAt())
	if err != nil {
		return nil, err
	}

	notifyAt, err := getTimeByTimestamp(eventIn.GetNotifyAt())
	if err != nil {
		return nil, err
	}

	var event *repository.Event
	if eventID == 0 {
		event, err = s.eUseCase.CreateEvent(ctx, eventIn.GetTitle(), eventIn.GetDescription(), startAt, endAt, notifyAt, userID)
	} else {
		event, err = s.eUseCase.UpdateEvent(ctx, eventIn.GetTitle(), eventIn.GetDescription(), startAt, endAt, notifyAt, eventID, userID)
	}

	if err != nil {
		oErr, ok := err.(outErr.OutError)
		if !ok {
			oErr = outErr.ErrInternal
			log.Errorf("unknown error: %v", err)
		}
		resp := &ModifyEventResponse{
			Result: &ModifyEventResponse_Error{
				Error: oErr.Error(),
			},
		}

		return resp, nil
	}
	outEvent, err := getGrpcEventByRepoEvent(event)
	if err != nil {
		return nil, err
	}

	resp := &ModifyEventResponse{
		Result: &ModifyEventResponse_Event{
			Event: outEvent,
		},
	}

	return resp, nil
}

func (s *CalendarServer) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteResponse, error) {
	err := s.eUseCase.DeleteEvent(ctx, int(req.GetEventId()))
	if err != nil {
		oErr, ok := err.(outErr.OutError)
		if !ok {
			oErr = outErr.ErrInternal
			log.Errorf("unknown error: %v", err)
		}
		resp := &DeleteResponse{
			Result: &DeleteResponse_Error{
				Error: oErr.Error(),
			},
		}

		return resp, nil
	}

	resp := &DeleteResponse{
		Result: &DeleteResponse_Success{
			Success: true,
		},
	}

	return resp, nil
}

func (s *CalendarServer) getFilteredEvents(ctx context.Context, req *GetEventsRequest, f filterFunc) (*GetEventsResponse, error) {
	userID, err := getUser(ctx)
	if err != nil {
		return nil, err
	}

	beginAt, err := getTimeByTimestamp(req.GetBeginAt())
	if err != nil {
		return nil, err
	}

	events, err := f(ctx, userID, beginAt)
	if err != nil {
		oErr, ok := err.(outErr.OutError)
		if !ok {
			oErr = outErr.ErrInternal
			log.Errorf("unknown error: %v", err)
		}
		resp := &GetEventsResponse{
			Result: &GetEventsResponse_Error{
				Error: oErr.Error(),
			},
		}

		return resp, nil
	}

	outEvents := new(EventList)
	for _, event := range events {
		outEvent, err := getGrpcEventByRepoEvent(event)
		if err != nil {
			return nil, err
		}
		outEvents.Events = append(outEvents.Events, outEvent)
	}

	resp := &GetEventsResponse{
		Result: &GetEventsResponse_Events{
			Events: outEvents,
		},
	}

	return resp, nil
}

func (s *CalendarServer) GetEventsDay(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error) {
	return s.getFilteredEvents(ctx, req, s.eUseCase.GetEventsDay)
}

func (s *CalendarServer) GetEventsWeek(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error) {
	return s.getFilteredEvents(ctx, req, s.eUseCase.GetEventsWeek)
}

func (s *CalendarServer) GetEventsMonth(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error) {
	return s.getFilteredEvents(ctx, req, s.eUseCase.GetEventsMonth)
}
