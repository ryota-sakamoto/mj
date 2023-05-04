package model

import (
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type ServerEvent interface {
	Into() *pb.RoomServerEvent
}

type ServerEventJoined struct {
	username string
}

func NewServerEventJoined(username string) ServerEvent {
	return &ServerEventJoined{
		username: username,
	}
}

func (s *ServerEventJoined) Into() *pb.RoomServerEvent {
	return &pb.RoomServerEvent{
		Event: &pb.RoomServerEvent_Joined{
			Joined: &pb.Joined{
				UserName: s.username,
			},
		},
	}
}

type ServerEventEmpty struct {
}

func NewServerEventEmpty() ServerEvent {
	return &ServerEventEmpty{}
}

func (s *ServerEventEmpty) Into() *pb.RoomServerEvent {
	return nil
}
