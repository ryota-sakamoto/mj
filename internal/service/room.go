package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ryota-sakamoto/mj/pkg/model"
)

type RoomService interface {
	Create(context.Context, *model.CreateRoom) (*model.Room, error)
}

type roomService struct {
}

func NewRoomService() RoomService {
	return &roomService{}
}

func (r *roomService) Create(ctx context.Context, req *model.CreateRoom) (*model.Room, error) {
	return &model.Room{
		ID:   uuid.NewString(),
		Name: req.Name,
	}, nil
}
