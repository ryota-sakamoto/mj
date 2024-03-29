package service

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/exp/slog"

	"github.com/ryota-sakamoto/mj/internal/repository"
	"github.com/ryota-sakamoto/mj/pkg/model"
)

type RoomService interface {
	Create(context.Context, *model.CreateRoom) (*model.Room, error)
	HandleUserEvent(context.Context, string, model.UserEvent) (model.ServerEvent, error)
	StreamServerEvent(context.Context, string) (model.ServerEvent, error)
}

type roomService struct {
	repository repository.RoomRepository
}

func NewRoomService(repository repository.RoomRepository) RoomService {
	return &roomService{
		repository: repository,
	}
}

func (r *roomService) Create(ctx context.Context, req *model.CreateRoom) (*model.Room, error) {
	return r.repository.Create(ctx, req)
}

func (r *roomService) HandleUserEvent(ctx context.Context, id string, event model.UserEvent) (model.ServerEvent, error) {
	slog.InfoCtx(ctx, "handle event",
		slog.Any("id", id),
		slog.Any("event", event),
	)

	switch e := event.(type) {
	case *model.UserEventJoin:
		return r.handleJoin(ctx, e)
	default:
		slog.ErrorCtx(ctx, "unknown event", slog.Any("event", event))

		return nil, fmt.Errorf("unknown event")
	}
}

func (r *roomService) handleJoin(ctx context.Context, req *model.UserEventJoin) (model.ServerEvent, error) {
	slog.InfoCtx(ctx, "receive join", slog.Any("req", req))

	room, err := r.repository.Get(ctx, req.ID, req.Password)
	if err != nil {
		return nil, fmt.Errorf("get room error: %w", err)
	}

	if err := r.repository.Join(ctx, room.ID, req.Username); err != nil {
		return nil, fmt.Errorf("join room error: %w", err)
	}

	return model.NewServerEventJoined(req.Username), nil
}

func (r *roomService) StreamServerEvent(ctx context.Context, id string) (model.ServerEvent, error) {
	select {
	case <-ctx.Done():
		return nil, io.EOF
	default:
	}

	return model.NewServerEventEmpty(), io.EOF
}
