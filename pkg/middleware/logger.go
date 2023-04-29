package middleware

import (
	"context"
	"time"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func Logger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		res, err := handler(ctx, req)
		if err != nil {
			slog.Error(
				"handler error",
				slog.Any("request_id", ctx.Value(requestIDKey{})),
				slog.Any("error", err),
			)
		}

		slog.Info(
			"access log",
			slog.Any("request_id", ctx.Value(requestIDKey{})),
			slog.String("method", info.FullMethod),
			slog.String("elapsed", time.Since(start).String()),
		)

		return res, err
	}
}
