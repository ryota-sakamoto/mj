package model

import (
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type Room struct {
	ID   string
	Name string
}

func (r *Room) Into() *pb.Room {
	return &pb.Room{
		Id:   r.ID,
		Name: r.Name,
	}
}

type CreateRoom struct {
	Name     string
	Password string
}

func FromCreateRoomRequest(r *pb.CreateRoomRequest) *CreateRoom {
	return &CreateRoom{
		Name:     r.Name,
		Password: r.Password,
	}
}
