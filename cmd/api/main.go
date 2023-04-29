package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ryota-sakamoto/mj/pkg/middleware"
	"github.com/ryota-sakamoto/mj/pkg/pb"
	"github.com/ryota-sakamoto/mj/pkg/service"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)))
}

func main() {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestID(),
			middleware.Logger(),
		),
	)
	reflection.Register(server)
	pb.RegisterRoomServiceServer(server, service.NewRoomService())

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go func() {
		slog.Info("starting server")
		if err := server.Serve(l); err != nil {
			panic(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	slog.Info("shutting down server")
	server.GracefulStop()
}
