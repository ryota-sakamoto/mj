package model

import (
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

type ServerEvent interface {
	Into() *pb.RoomServerEvent
}

type serverEventJoined struct {
	username string
}

func NewServerEventJoined(username string) ServerEvent {
	return &serverEventJoined{
		username: username,
	}
}

func (s *serverEventJoined) Into() *pb.RoomServerEvent {
	return &pb.RoomServerEvent{
		Event: &pb.RoomServerEvent_Joined{
			Joined: &pb.Joined{
				UserName: s.username,
			},
		},
	}
}

type serverEventEmpty struct {
}

func NewServerEventEmpty() ServerEvent {
	return &serverEventEmpty{}
}

func (s *serverEventEmpty) Into() *pb.RoomServerEvent {
	return nil
}

type serverEventRejected struct {
	reason string
}

func NewServerEventRejected(reason string) ServerEvent {
	return &serverEventRejected{
		reason: reason,
	}
}

func (s *serverEventRejected) Into() *pb.RoomServerEvent {
	return &pb.RoomServerEvent{
		Event: &pb.RoomServerEvent_Rejected{
			Rejected: &pb.Rejected{
				Reason: s.reason,
			},
		},
	}
}
