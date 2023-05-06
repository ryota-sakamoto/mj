package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/ryota-sakamoto/mj/pkg/model"
)

type RoomRepository interface {
	Create(context.Context, *model.CreateRoom) (*model.Room, error)
	Get(context.Context, string, string) (*model.Room, error)
	Join(context.Context, *model.Room, string) error
}

type roomRepository struct {
	sync.RWMutex

	rooms map[string]*innerRoom
}

func NewRoomRepository() RoomRepository {
	return &roomRepository{
		rooms: map[string]*innerRoom{},
	}
}

type innerRoom struct {
	room *model.Room
	salt string
	hash string
}

func hashWithSalt(password, salt string) string {
	return fmt.Sprintf("%+x", sha256.Sum256([]byte(password+salt)))
}

func newInnerRoom(req *model.CreateRoom) *innerRoom {
	id := uuid.NewString()
	salt := uuid.NewString()

	return &innerRoom{
		room: &model.Room{
			ID: id,
		},
		salt: salt,
		hash: hashWithSalt(req.Password, salt),
	}
}

func (r *innerRoom) match(password string) bool {
	return hashWithSalt(password, r.salt) == r.hash
}

func (r *roomRepository) Create(ctx context.Context, req *model.CreateRoom) (*model.Room, error) {
	r.Lock()
	defer r.Unlock()

	inner := newInnerRoom(req)
	r.rooms[inner.room.ID] = inner

	return inner.room, nil
}

func (r *roomRepository) Get(ctx context.Context, id string, password string) (*model.Room, error) {
	r.RLock()
	defer r.RUnlock()

	room, ok := r.rooms[id]
	if !ok || !room.match(password) {
		return nil, model.ErrNotFound
	}

	return room.room, nil
}

func (r *roomRepository) Join(ctx context.Context, room *model.Room, username string) error {
	r.Lock()
	defer r.Unlock()

	return nil
}
