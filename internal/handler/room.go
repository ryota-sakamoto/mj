package handler

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/ryota-sakamoto/mj/internal/service"
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type RoomHandler struct {
	pb.UnimplementedRoomServiceServer

	service service.RoomService
}

func NewRoomHandler(service service.RoomService) *RoomHandler {
	return &RoomHandler{
		service: service,
	}
}

func (r *RoomHandler) Create(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	slog.Info("create", slog.Any("req", req))

	return &pb.CreateRoomResponse{}, nil
}
