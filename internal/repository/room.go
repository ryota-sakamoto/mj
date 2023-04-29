package repository

import (
	"context"
	"sync"

	"github.com/ryota-sakamoto/mj/pkg/model"
)

type RoomRepository interface {
	Create(context.Context, *model.CreateRoom) error
}

type roomRepository struct {
	sync.Mutex
}

func NewRoomRepository() RoomRepository {
	return &roomRepository{}
}

func (r *roomRepository) Create(ctx context.Context, room *model.CreateRoom) error {
	r.Lock()
	defer r.Unlock()

	return nil
}
