package manager

import (
	"context"
	"io"

	"golang.org/x/exp/slog"

	"github.com/ryota-sakamoto/mj/internal/service"
	"github.com/ryota-sakamoto/mj/pkg/model"
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type UserManager interface {
	HandleEvent()
}

type userManager struct {
	roomID string
	userID string

	stream  pb.RoomService_StreamEventsServer
	service service.RoomService
}

func NewUserManager(stream pb.RoomService_StreamEventsServer, service service.RoomService) UserManager {
	return &userManager{
		stream:  stream,
		service: service,
	}
}

func (u *userManager) HandleEvent() {
	for {
		event, err := u.stream.Recv()
		if err != nil {
			return
		}

		joinEvent := event.GetJoin()
		if joinEvent == nil {
			continue
		}

		u.roomID = joinEvent.Id
		u.userID = joinEvent.UserName

		if err := u.handleRoomUserEvent(u.stream.Context(), event); err != nil {
			continue
		}

		break
	}

	go u.handleUserEvent()
	go u.streamServerEvent()

	<-u.stream.Context().Done()
}

func (u *userManager) handleUserEvent() {
	ctx := u.stream.Context()

	for {
		e, err := u.stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			slog.ErrorCtx(ctx, "stream user event error", slog.Any("error", err))
			break
		}

		u.handleRoomUserEvent(ctx, e)
	}
}

func (u *userManager) handleRoomUserEvent(ctx context.Context, e *pb.RoomUserEvent) error {
	event, err := model.FromUserEvent(e)
	if err != nil {
		slog.ErrorCtx(ctx, "parse user event error", slog.Any("error", err))
		return err
	}

	result, err := u.service.HandleUserEvent(ctx, u.roomID, event)
	if err != nil {
		slog.ErrorCtx(ctx, "handle user event error", slog.Any("error", err))
		return err
	}

	slog.InfoCtx(ctx, "handle user event result", slog.Any("result", result.Into()))
	u.stream.Send(result.Into())

	return nil
}

func (u *userManager) streamServerEvent() {
	ctx := u.stream.Context()

	for {
		event, err := u.service.StreamServerEvent(ctx, u.roomID)
		if err != nil {
			if err == io.EOF {
				break
			}

			slog.ErrorCtx(ctx, "stream server event error", slog.Any("error", err))
			break
		}

		slog.InfoCtx(ctx, "handle server event result", slog.Any("event", event))
		u.stream.Send(event.Into())
	}
}
