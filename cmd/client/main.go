package main

import (
	"context"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ryota-sakamoto/mj/pkg/pb"
)

func main() {
	conn, err := grpc.Dial(
		":8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	client := pb.NewRoomServiceClient(conn)
	room, err := client.Create(context.TODO(), &pb.CreateRoomRequest{
		Password:    "test",
		OwnerName:   "owner",
		PlayerCount: 4,
	})
	if err != nil {
		panic(err)
	}

	slog.Info("create room res", slog.Any("room", room))

	eventClient, err := client.StreamEvents(context.Background())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			event, err := eventClient.Recv()
			if err != nil {
				panic(err)
			}

			slog.Info("receive event", slog.Any("event", event))
		}
	}()

	for {
		err = eventClient.Send(&pb.RoomUserEvent{
			Event: &pb.RoomUserEvent_Join{
				Join: &pb.Join{
					Id:       room.Id,
					UserName: "user",
					Password: "test",
				},
			},
		})
		if err != nil {
			panic(err)
		}
	}

	eventClient.CloseSend()
	<-eventClient.Context().Done()
}
