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
			return nil, transformError(err)
		}

		return res, nil
	}
}

func TransformStreamErrorResponse() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err != nil {
			return transformError(err)
		}

		return nil
	}
}

func transformError(err error) error {
	if errors.Is(err, model.ErrNotFound) {
		st := status.New(codes.NotFound, err.Error())
		return st.Err()
	}

	if e, ok := err.(model.ValidationError); ok {
		st := status.New(codes.InvalidArgument, e.Error())
		return st.Err()
	}

	return err
}
