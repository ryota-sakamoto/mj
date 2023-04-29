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
