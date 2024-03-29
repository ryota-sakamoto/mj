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

	"github.com/ryota-sakamoto/mj/internal/handler"
	"github.com/ryota-sakamoto/mj/internal/repository"
	"github.com/ryota-sakamoto/mj/internal/service"
	"github.com/ryota-sakamoto/mj/pkg/middleware"
	"github.com/ryota-sakamoto/mj/pkg/pb"
)

func init() {
	slog.SetDefault(slog.New(middleware.NewLogHandler()))
}

func main() {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Logger(),
			middleware.TransformErrorResponse(),
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamLogger(),
			middleware.TransformStreamErrorResponse(),
		),
	)
	reflection.Register(server)
	pb.RegisterRoomServiceServer(server, handler.NewRoomHandler(service.NewRoomService(repository.NewRoomRepository())))

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
