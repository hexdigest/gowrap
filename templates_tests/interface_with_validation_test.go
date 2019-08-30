package templatestests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// ValidatableRequest has Validate() error method that is called by the "validate" decorator
type ValidatableRequest struct {
	err error
}

func (v *ValidatableRequest) Validate() error {
	return v.err
}

//TestInterface is used to test templates
type InterfaceWithValidtableArg interface {
	Method(ctx context.Context, r *ValidatableRequest) error
}

func (v *ValidatableRequest) Method(context.Context, *ValidatableRequest) error {
	return nil
}

func TestInterfaceWithValidatableArg_Method(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := ValidatableRequest{}
		wrapped := NewInterfaceWithValidtableArgWithValidation(&r)
		err := wrapped.Method(context.Background(), &r)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		r := ValidatableRequest{
			err: errors.New("unexpected error"),
		}
		wrapped := NewInterfaceWithValidtableArgWithValidation(&r)
		err := wrapped.Method(context.Background(), &r)
		require.Error(t, err)
	})
}
