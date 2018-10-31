package templatestests

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testImpl struct {
	delay       time.Duration
	ch          chan struct{}
	err         error
	r1, r2      string
	callCounter uint64
}

func (f *testImpl) F(ctx context.Context, a1 string, a2 ...string) (r1, r2 string, err error) {
	atomic.AddUint64(&f.callCounter, 1)

	if f.ch != nil {
		defer close(f.ch)
	}

	err = f.err
	r1 = f.r1
	r2 = f.r2

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-time.After(f.delay):
	}

	return
}

func (f *testImpl) NoError(s string) string {
	return s
}

func (f *testImpl) NoParamsOrResults() {}

func TestTestInterfaceWithFallback_F(t *testing.T) {
	t.Run("one implementation success", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2"}
		wrapped := NewTestInterfaceWithFallback(time.Second, impl)

		r1, r2, err := wrapped.F(context.Background(), "")
		require.NoError(t, err, "%T", err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)
	})

	t.Run("one implementation failure", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2", err: errors.New("failure")}
		wrapped := NewTestInterfaceWithFallback(time.Second, impl)

		_, _, err := wrapped.F(context.Background(), "")
		assert.Error(t, err)
	})

	t.Run("first comes first, no errors", func(t *testing.T) {
		impl := &testImpl{r1: "11", r2: "12"}
		impl2 := &testImpl{r1: "21", r2: "22", delay: time.Second}
		wrapped := NewTestInterfaceWithFallback(time.Second, impl, impl2)

		r1, r2, err := wrapped.F(context.Background(), "")
		require.NoError(t, err)
		assert.Equal(t, "11", r1)
		assert.Equal(t, "12", r2)
	})

	t.Run("second comes first, no errors", func(t *testing.T) {
		impl := &testImpl{r1: "11", r2: "12", delay: time.Second}
		impl2 := &testImpl{r1: "21", r2: "22"}
		wrapped := NewTestInterfaceWithFallback(100*time.Millisecond, impl, impl2)

		r1, r2, err := wrapped.F(context.Background(), "")
		require.NoError(t, err)
		assert.Equal(t, "21", r1)
		assert.Equal(t, "22", r2)
	})

	t.Run("first quickly returns error", func(t *testing.T) {
		impl := &testImpl{r1: "11", r2: "12", err: errors.New("failure")}
		impl2 := &testImpl{r1: "21", r2: "22", delay: 50 * time.Millisecond}
		wrapped := NewTestInterfaceWithFallback(100*time.Millisecond, impl, impl2)

		r1, r2, err := wrapped.F(context.Background(), "")
		require.NoError(t, err)
		assert.Equal(t, "21", r1)
		assert.Equal(t, "22", r2)
	})

	t.Run("first quickly returns error", func(t *testing.T) {
		impl := &testImpl{r1: "11", r2: "12", delay: 20 * time.Millisecond, err: errors.New("first failure")}
		impl2 := &testImpl{r1: "21", r2: "22", delay: 10 * time.Millisecond, err: errors.New("second failure")}
		wrapped := NewTestInterfaceWithFallback(100*time.Millisecond, impl, impl2)

		_, _, err := wrapped.F(context.Background(), "")
		require.Error(t, err)
		assert.Equal(t, "*templatestests.testImpl: first failure;*templatestests.testImpl: second failure", err.Error())
	})
}
