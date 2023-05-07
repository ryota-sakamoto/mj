package model

import (
	"fmt"

	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type Room struct {
	ID          string
	PlayerCount int64
}

func (r *Room) Into() *pb.Room {
	return &pb.Room{
		Id:          r.ID,
		PlayerCount: r.PlayerCount,
	}
}

type CreateRoom struct {
	Password    string
	OwnerName   string
	PlayerCount int64
}

func FromCreateRoomRequest(r *pb.CreateRoomRequest) (*CreateRoom, error) {
	if r.PlayerCount < 1 || r.PlayerCount > 4 {
		return nil, NewValidationError(fmt.Sprintf("player count is not within the range from 1 to 4(value: %d)", r.PlayerCount))
	}

	return &CreateRoom{
		Password:    r.Password,
		OwnerName:   r.OwnerName,
		PlayerCount: r.PlayerCount,
	}, nil
}
