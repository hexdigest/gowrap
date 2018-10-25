package templatestests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestInterfaceWithRetry_F(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2"}
		wrapped := NewTestInterfaceWithRetry(impl, 2, time.Second)

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")
		assert.NoError(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)
	})

	t.Run("error", func(t *testing.T) {
		errUnexpected := errors.New("unexpected error")
		impl := &testImpl{r1: "1", r2: "2", err: errUnexpected}
		wrapped := NewTestInterfaceWithRetry(impl, 1, time.Millisecond)

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")
		require.Error(t, err)
		assert.Equal(t, errUnexpected, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)
		assert.EqualValues(t, 2, impl.callCounter)
	})

	t.Run("error and context deadline", func(t *testing.T) {
		errUnexpected := errors.New("unexpected error")
		impl := &testImpl{r1: "1", r2: "2", err: errUnexpected}
		wrapped := NewTestInterfaceWithRetry(impl, 1, time.Second)

		ctx, cancelFunc := context.WithCancel(context.Background())
		cancelFunc()

		r1, r2, err := wrapped.F(ctx, "a1", "a2")
		require.Error(t, err)
		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)
		assert.EqualValues(t, 1, impl.callCounter)
	})
}
