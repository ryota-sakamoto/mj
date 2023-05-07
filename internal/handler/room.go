package handler

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/ryota-sakamoto/mj/internal/manager"
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

	innerReq, err := model.FromCreateRoomRequest(req)
	if err != nil {
		return nil, err
	}

	res, err := r.service.Create(ctx, innerReq)
	if err != nil {
		return nil, err
	}

	return res.Into(), nil
}

func (r *RoomHandler) StreamEvents(s pb.RoomService_StreamEventsServer) error {
	m := manager.NewUserManager(s, r.service)
	go m.HandleEvent()

	<-s.Context().Done()
	slog.InfoCtx(s.Context(), "close stream events")

	return nil
}
