package model

import (
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type Room struct {
	ID string
}

func (r *Room) Into() *pb.Room {
	return &pb.Room{
		Id: r.ID,
	}
}

type CreateRoom struct {
	Password string
}

func FromCreateRoomRequest(r *pb.CreateRoomRequest) *CreateRoom {
	return &CreateRoom{
		Password: r.Password,
	}
}

type JoinRoom struct {
	ID       string
	Password string
	UserName string
}

func FromJoinRoomRequest(r *pb.JoinRoomRequest) *JoinRoom {
	return &JoinRoom{
		ID:       r.Id,
		Password: r.Password,
		UserName: r.UserName,
	}
}
