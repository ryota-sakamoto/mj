package service

import (
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type RoomService struct {
	pb.UnimplementedRoomServiceServer
}

func NewRoomService() *RoomService {
	return &RoomService{}
}
