package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type requestIDKey struct{}

func RequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = context.WithValue(ctx, requestIDKey{}, uuid.NewString())

		return handler(ctx, req)
	}
}
