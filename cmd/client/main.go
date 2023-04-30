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
		Password:  "test",
		OwnerName: "owner",
	})
	if err != nil {
		panic(err)
	}

	slog.Info("create room res", slog.Any("room", room))

	res, err := client.Join(context.TODO(), &pb.JoinRoomRequest{
		Id:       room.Id,
		Password: "test",
		UserName: "user",
	})
	if err != nil {
		panic(err)
	}

	slog.Info("join room res", slog.Any("res", res))
}
