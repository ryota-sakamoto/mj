package middleware

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ryota-sakamoto/mj/pkg/model"
)

func TransformErrorResponse() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		res, err := handler(ctx, req)
		if err != nil {
			return nil, TransformError(err)
		}

		return res, nil
	}
}

func TransformStreamErrorResponse() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err != nil {
			return TransformError(err)
		}

		return nil
	}
}

func TransformError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, model.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if e, ok := err.(model.ValidationError); ok {
		return status.Error(codes.InvalidArgument, e.Error())
	}

	return status.Error(codes.Internal, "internal error")
}
