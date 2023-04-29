package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ryota-sakamoto/mj/internal/repository"
	"github.com/ryota-sakamoto/mj/pkg/model"
)

type RoomService interface {
	Create(context.Context, *model.CreateRoom) (*model.Room, error)
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
	if err := r.repository.Create(ctx, req); err != nil {
		return nil, err
	}

	return &model.Room{
		ID:   uuid.NewString(),
		Name: req.Name,
	}, nil
}
