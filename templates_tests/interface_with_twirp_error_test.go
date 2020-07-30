package templatestests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/twitchtv/twirp"
	"testing"
)

type WithTwirpError struct{}

type MethodRequest struct {
	Foo string `json:"foo"`
}

type MethodResponse struct{}

//TestInterface is used to test templates
type InterfaceWithTwirpError interface {
	Method(ctx context.Context, r *MethodRequest) (*MethodResponse, error)
}

func (v *WithTwirpError) Method(_ context.Context, req *MethodRequest) (*MethodResponse, error) {
	if req.Foo != "bar" {
		return nil, fmt.Errorf("foo != bar")
	}
	return nil, nil
}

func TestInterfaceWithTwirpError_Method(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		wrapped := NewInterfaceWithTwirpErrorWithTwirpError(&WithTwirpError{})
		_, err := wrapped.Method(context.Background(), &MethodRequest{Foo: "bar"})
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		wrapped := NewInterfaceWithTwirpErrorWithTwirpError(&WithTwirpError{})
		_, err := wrapped.Method(context.Background(), &MethodRequest{Foo: "invalid"})
		require.Error(t, err)
		twirpErr, ok := err.(twirp.Error)
		require.True(t, ok)
		require.Equal(t, twirp.Internal, twirpErr.Code())
		require.Equal(t, "{\"foo\":\"invalid\"}", twirpErr.Meta("request"))
	})
}
