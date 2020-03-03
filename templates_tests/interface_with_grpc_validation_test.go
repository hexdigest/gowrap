package templatestests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestInterfaceWithValidatableArgWithGRPC_Method(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := ValidatableRequest{}
		wrapped := NewInterfaceWithValidtableArgWithGRPCValidation(&r)
		err := wrapped.Method(context.Background(), &r)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		r := ValidatableRequest{
			err: errors.New("unexpected error"),
		}
		wrapped := NewInterfaceWithValidtableArgWithGRPCValidation(&r)
		err := wrapped.Method(context.Background(), &r)
		require.Error(t, err)
		s, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, s.Code())
	})
}
