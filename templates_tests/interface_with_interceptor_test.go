package templatestests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestInterfaceWithInterceptor_F(t *testing.T) {
	i := NewInterceptorTestInterface(&testImpl{})

	called := false
	defer func() {
		require.True(t, called)
	}()

	i.InterceptorFuncF = func(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
		called = true
		return "", "", nil // or we can do : return m._base.F(ctx, a1, a2...)
	}
	i.F(context.Background(), "", "")
}
