package handler

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/ryota-sakamoto/mj/internal/service"
	"github.com/ryota-sakamoto/mj/pkg/model"
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

func (r *RoomHandler) Create(ctx context.Context, req *pb.CreateRoomRequest) (*pb.Room, error) {
	slog.InfoCtx(ctx, "create request", slog.Any("req", req))

	res, err := r.service.Create(ctx, model.FromCreateRoomRequest(req))
	if err != nil {
		return nil, err
	}

	return res.Into(), nil
}

func (r *RoomHandler) Join(ctx context.Context, req *pb.JoinRoomRequest) (*pb.JoinRoom, error) {
	err := r.service.Join(ctx, model.FromJoinRoomRequest(req))
	if err != nil {
		return nil, err
	}

	return &pb.JoinRoom{
		Token: "hoge",
	}, nil
}
