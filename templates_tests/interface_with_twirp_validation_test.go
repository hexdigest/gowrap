package templatestests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twitchtv/twirp"
)

func TestInterfaceWithValidatableArgWithTwirp_Method(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := ValidatableRequest{}
		wrapped := NewInterfaceWithValidtableArgWithTwirpValidation(&r)
		err := wrapped.Method(context.Background(), &r)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		r := ValidatableRequest{
			err: errors.New("unexpected error"),
		}
		wrapped := NewInterfaceWithValidtableArgWithTwirpValidation(&r)
		err := wrapped.Method(context.Background(), &r)
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		require.Equal(t, twirp.InvalidArgument, twirpErr.Code())
	})
}
