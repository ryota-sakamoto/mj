package handler

import (
	"context"
	"errors"
	"io"

	"golang.org/x/exp/slog"

	"github.com/ryota-sakamoto/mj/internal/service"
	"github.com/ryota-sakamoto/mj/pkg/model"
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type RoomHandler struct {
	pb.UnimplementedRoomServiceServer

	service service.RoomService
}

func NewRoomHandler(service service.RoomService) *RoomHandler {
	return &RoomHandler{
		service: service,
	}
}

func (r *RoomHandler) Create(ctx context.Context, req *pb.CreateRoomRequest) (*pb.Room, error) {
	slog.InfoCtx(ctx, "create request", slog.Any("req", req))

	innerReq, err := model.FromCreateRoomRequest(req)
	if err != nil {
		return nil, err
	}

	res, err := r.service.Create(ctx, innerReq)
	if err != nil {
		return nil, err
	}

	return res.Into(), nil
}

func (r *RoomHandler) StreamEvents(s pb.RoomService_StreamEventsServer) error {
	e, err := s.Recv()
	if err != nil {
		return err
	}

	result, id, err := r.initStream(s.Context(), e)
	if err != nil {
		return err
	}

	s.Send(result.Into())

	go r.handleUserEvent(s, id)
	go r.streamServerEvent(s, id)

	<-s.Context().Done()
	slog.InfoCtx(s.Context(), "close stream events")

	return nil
}

func (r *RoomHandler) initStream(ctx context.Context, event *pb.RoomUserEvent) (model.ServerEvent, string, error) {
	joinEvent := event.GetJoin()
	if joinEvent == nil {
		slog.ErrorCtx(ctx, "the request is not joined to the room")

		return nil, "", errors.New("the request is not joined to the room")
	}

	e, err := model.NewUserEventJoin(joinEvent)
	if joinEvent == nil {
		return nil, "", err
	}

	result, err := r.service.HandleUserEvent(ctx, e.ID, e)
	if err != nil {
		return nil, "", err
	}

	return result, e.ID, nil
}

func (r *RoomHandler) handleUserEvent(s pb.RoomService_StreamEventsServer, id string) {
	ctx := s.Context()

	for {
		e, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			slog.ErrorCtx(ctx, "stream user event error", slog.Any("error", err))
			break
		}

		event, err := model.FromUserEvent(e)
		if err != nil {
			slog.ErrorCtx(ctx, "parse user event error", slog.Any("error", err))
			continue
		}

		result, err := r.service.HandleUserEvent(ctx, id, event)
		if err != nil {
			slog.ErrorCtx(ctx, "handle user event error", slog.Any("error", err))
			continue
		}

		slog.InfoCtx(ctx, "handle user event result", slog.Any("result", result.Into()))
		s.Send(result.Into())
	}
}

func (r *RoomHandler) streamServerEvent(s pb.RoomService_StreamEventsServer, id string) {
	ctx := s.Context()

	for {
		event, err := r.service.StreamServerEvent(ctx, id)
		if err != nil {
			if err == io.EOF {
				break
			}

			slog.ErrorCtx(ctx, "stream server event error", slog.Any("error", err))
			break
		}

		slog.InfoCtx(ctx, "handle server event result", slog.Any("event", event))
		s.Send(event.Into())
	}
}
