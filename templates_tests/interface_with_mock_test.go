package templatestests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestInterfaceWithMock_F(t *testing.T) {
	m := NewMockTestInterface()

	called := false
	defer func() {
		require.True(t, called)
		require.Equal(t, 1, m.NCalledMockF)
	}()

	m.MockFuncF = func(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
		called = true
		return "", "", nil
	}
	m.F(context.Background(), "", "")
}
