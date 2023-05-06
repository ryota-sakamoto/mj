package middleware

import (
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ryota-sakamoto/mj/pkg/model"
)

func TransformStreamResponse() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				st := status.New(codes.NotFound, err.Error())
				return st.Err()
			}

			return err
		}

		return nil
	}
}
