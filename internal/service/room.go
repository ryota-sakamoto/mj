package service

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/ryota-sakamoto/mj/internal/repository"
	"github.com/ryota-sakamoto/mj/pkg/model"
)

type RoomService interface {
	Create(context.Context, *model.CreateRoom) (*model.Room, error)
	Join(context.Context, *model.JoinRoom) error
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

func (r *roomService) Join(ctx context.Context, req *model.JoinRoom) error {
	_, err := r.repository.Get(ctx, req.ID, req.Password)
	if err != nil {
		slog.ErrorCtx(ctx,
			"get room error",
			slog.String("id", req.ID),
			slog.Any("error", err),
		)

		return err
	}

	return nil
}
