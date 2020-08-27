package grpc

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/server"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type filterFunc func(context.Context, int, time.Time) ([]*repository.Event, error)

type CalendarServer struct {
	gs       *grpc.Server
	eUseCase usecase.EventUseCaser
}

func NewCalendarServer(eUseCase usecase.EventUseCaser) *CalendarServer {
	s := new(CalendarServer)
	s.eUseCase = eUseCase

	return s
}

func (s *CalendarServer) Run(addr string) error {
	s.gs = grpc.NewServer(unaryInterceptor())
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	RegisterCalendarServiceServer(s.gs, s)

	return s.gs.Serve(l)
}

func (s *CalendarServer) Down() {
	s.gs.GracefulStop()
}

func unaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		defer log.Infof(
			"%s %v",
			info.FullMethod,
			time.Since(start),
		)

		resp, err := handler(ctx, req)
		if err != nil {
			log.Errorf("method %q throws error: %v", info.FullMethod, err)
		}

		return resp, err
	})
}

func getUser(ctx context.Context) (int, error) {
	var userID int

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return userID, fmt.Errorf("unable to read incoming context")
	}
	if uID := meta.Get(server.UserIDHeaderKey); uID != nil {
		var err error
		userID, err = strconv.Atoi(uID[0])
		if err != nil {
			return userID, err
		}
	}
	if userID == 0 {
		return userID, fmt.Errorf("invalid user")
	}

	return userID, nil
}

func getTimeByTimestamp(t *timestamp.Timestamp) (time.Time, error) {
	var timeAt time.Time
	var err error
	if t != nil {
		timeAt, err = ptypes.Timestamp(t)
		if err != nil {
			return timeAt, err
		}
	}

	return timeAt, nil
}

func getGrpcEventByRepoEvent(event *repository.Event) (*Event, error) {
	var err error
	outEvent := &Event{
		Id:          int32(event.ID),
		Title:       event.Title,
		Description: event.Description,
		UserId:      int32(event.UserID),
	}
	if outEvent.StartAt, err = ptypes.TimestampProto(event.StartAt); err != nil {
		return nil, err
	}
	if outEvent.EndAt, err = ptypes.TimestampProto(event.EndAt); err != nil {
		return nil, err
	}
	if outEvent.NotifyAt, err = ptypes.TimestampProto(event.NotifyAt); err != nil {
		return nil, err
	}

	return outEvent, nil
}
