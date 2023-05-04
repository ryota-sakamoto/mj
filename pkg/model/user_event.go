package model

import (
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type UserEvent interface {
	userEventImpl
}

type userEventImpl interface {
	userEvent()
}

func FromUserEvent(event *pb.RoomUserEvent) (UserEvent, error) {
	switch event.Event.(type) {
	case *pb.RoomUserEvent_Join:
		return NewUserEventJoin(event.GetJoin())
	default:
	}

	return nil, nil
}

type UserEventJoin struct {
	userEventImpl

	ID       string
	Password string
	Username string
}

func NewUserEventJoin(event *pb.Join) (*UserEventJoin, error) {
	return &UserEventJoin{
		ID:       event.Id,
		Password: event.Password,
		Username: event.UserName,
	}, nil
}
