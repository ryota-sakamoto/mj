package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/ryota-sakamoto/mj/pkg/model"
)

type RoomRepository interface {
	Create(context.Context, *model.CreateRoom) (*model.Room, error)
	Get(context.Context, string, string) (*model.Room, error)
	Join(context.Context, string, string) error
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
	room  *model.Room
	users []*model.User
	salt  string
	hash  string
}

func hashWithSalt(password, salt string) string {
	return fmt.Sprintf("%+x", sha256.Sum256([]byte(password+salt)))
}

func newInnerRoom(req *model.CreateRoom) *innerRoom {
	id := uuid.NewString()
	salt := uuid.NewString()

	return &innerRoom{
		room: &model.Room{
			ID:          id,
			PlayerCount: req.PlayerCount,
		},
		users: []*model.User{
			{
				Name: req.OwnerName,
			},
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

func (r *roomRepository) Join(ctx context.Context, id string, username string) error {
	r.Lock()
	defer r.Unlock()

	room := r.rooms[id]
	for _, u := range room.users {
		if u.Name == username {
			slog.InfoCtx(
				ctx,
				"user is alread joined",
				slog.String("id", id),
				slog.String("username", username),
			)
			return model.ErrAlreadyJoined
		}
	}

	if len(room.users) >= int(room.room.PlayerCount) {
		return model.ErrLimitExceeded
	}

	room.users = append(room.users, &model.User{
		Name: username,
	})

	r.rooms[id] = room

	return nil
}
