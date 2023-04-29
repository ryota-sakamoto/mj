package middleware

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

type requestIDKey struct{}

func Logger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		ctx = context.WithValue(ctx, requestIDKey{}, uuid.NewString())

		res, err := handler(ctx, req)
		if err != nil {
			slog.ErrorCtx(
				ctx,
				"handler error",
				slog.Any("error", err),
			)
		}

		slog.InfoCtx(
			ctx,
			"access log",
			slog.String("method", info.FullMethod),
			slog.String("elapsed", time.Since(start).String()),
		)

		return res, err
	}
}

type Handler struct {
	*slog.JSONHandler
}

func NewLogHandler() slog.Handler {
	return &Handler{
		JSONHandler: slog.NewJSONHandler(os.Stdout),
	}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, b := ctx.Value(requestIDKey{}).(string); b {
		r.Add(slog.Any("request_id", requestID))
	}

	return h.JSONHandler.Handle(ctx, r)
}
