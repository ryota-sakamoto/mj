package middleware_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ryota-sakamoto/mj/pkg/middleware"
	"github.com/ryota-sakamoto/mj/pkg/model"
)

func TestTransformError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "error is nil",
			err:      nil,
			expected: nil,
		},
		{
			name:     "not found error",
			err:      fmt.Errorf("hoge is not found: %w", model.ErrNotFound),
			expected: status.Errorf(codes.NotFound, "hoge is not found: %s", model.ErrNotFound.Error()),
		},
		{
			name:     "validation error",
			err:      model.NewValidationError("hoge is invalid"),
			expected: status.Error(codes.InvalidArgument, "hoge is invalid"),
		},
		{
			name:     "internal error",
			err:      io.EOF,
			expected: status.Error(codes.Internal, "internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := middleware.TransformError(tt.err)
			if tt.err == nil {
				assert.Equal(t, tt.expected, actual)
			} else {
				assert.Equal(t, tt.expected.Error(), actual.Error())
			}
		})
	}
}
