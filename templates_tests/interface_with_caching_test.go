package templatestests

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type cachingTestImpl struct {
	fCallCount          int32
	noErrorCallCount    int32
	contextNoErrorCount int32
	noParamsCallCount   int32
	channelsCallCount   int32
	shouldError         bool
}

var errCaching = errors.New("caching test error")

func (c *cachingTestImpl) F(ctx context.Context, a1 string, a2 ...string) (result1, result2 string, err error) {
	atomic.AddInt32(&c.fCallCount, 1)
	if c.shouldError {
		return "", "", errCaching
	}
	return "result1_" + a1, "result2_" + a1, nil
}

func (c *cachingTestImpl) NoError(s string) string {
	atomic.AddInt32(&c.noErrorCallCount, 1)
	return "noerror_" + s
}

func (c *cachingTestImpl) ContextNoError(ctx context.Context, a1 string, a2 string) {
	atomic.AddInt32(&c.contextNoErrorCount, 1)
}

func (c *cachingTestImpl) NoParamsOrResults() {
	atomic.AddInt32(&c.noParamsCallCount, 1)
}

func (c *cachingTestImpl) Channels(chA chan bool, chB chan<- bool, chanC <-chan bool) {
	atomic.AddInt32(&c.channelsCallCount, 1)
}

func TestTestInterfaceWithCaching_F(t *testing.T) {
	ctx := context.Background()

	t.Run("caches successful results", func(t *testing.T) {
		impl := &cachingTestImpl{}
		wrapped := NewTestInterfaceWithCaching(impl, time.Hour, time.Hour)

		// First call
		r1, r2, err := wrapped.F(ctx, "test")
		assert.NoError(t, err)
		assert.Equal(t, "result1_test", r1)
		assert.Equal(t, "result2_test", r2)
		assert.Equal(t, int32(1), atomic.LoadInt32(&impl.fCallCount))

		// Second call with same params - should be cached
		r1, r2, err = wrapped.F(ctx, "test")
		assert.NoError(t, err)
		assert.Equal(t, "result1_test", r1)
		assert.Equal(t, "result2_test", r2)
		assert.Equal(t, int32(1), atomic.LoadInt32(&impl.fCallCount)) // Still 1
	})

	t.Run("different parameters create different cache keys", func(t *testing.T) {
		impl := &cachingTestImpl{}
		wrapped := NewTestInterfaceWithCaching(impl, time.Hour, time.Hour)

		// Call with "test1"
		r1, r2, err := wrapped.F(ctx, "test1")
		assert.NoError(t, err)
		assert.Equal(t, "result1_test1", r1)
		assert.Equal(t, "result2_test1", r2)
		assert.Equal(t, int32(1), atomic.LoadInt32(&impl.fCallCount))

		// Call with "test2" - different param, should call impl
		r1, r2, err = wrapped.F(ctx, "test2")
		assert.NoError(t, err)
		assert.Equal(t, "result1_test2", r1)
		assert.Equal(t, "result2_test2", r2)
		assert.Equal(t, int32(2), atomic.LoadInt32(&impl.fCallCount))
	})

	t.Run("errors are not cached", func(t *testing.T) {
		impl := &cachingTestImpl{shouldError: true}
		wrapped := NewTestInterfaceWithCaching(impl, time.Hour, time.Hour)

		// First call - should error
		_, _, err := wrapped.F(ctx, "test")
		assert.Equal(t, errCaching, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&impl.fCallCount))

		// Second call - should call impl again (errors not cached)
		_, _, err = wrapped.F(ctx, "test")
		assert.Equal(t, errCaching, err)
		assert.Equal(t, int32(2), atomic.LoadInt32(&impl.fCallCount))
	})
}

func TestTestInterfaceWithCaching_NoError(t *testing.T) {
	t.Run("caches single return value methods", func(t *testing.T) {
		impl := &cachingTestImpl{}
		wrapped := NewTestInterfaceWithCaching(impl, time.Hour, time.Hour)

		// First call
		result := wrapped.NoError("test")
		assert.Equal(t, "noerror_test", result)
		assert.Equal(t, int32(1), atomic.LoadInt32(&impl.noErrorCallCount))

		// Second call - should be cached
		result = wrapped.NoError("test")
		assert.Equal(t, "noerror_test", result)
		assert.Equal(t, int32(1), atomic.LoadInt32(&impl.noErrorCallCount)) // Still 1
	})
}

func TestTestInterfaceWithCaching_NoReturns(t *testing.T) {
	t.Run("methods without return values are not cached", func(t *testing.T) {
		impl := &cachingTestImpl{}
		wrapped := NewTestInterfaceWithCaching(impl, time.Hour, time.Hour)

		// First call
		wrapped.ContextNoError(context.Background(), "a", "b")
		assert.Equal(t, int32(1), atomic.LoadInt32(&impl.contextNoErrorCount))

		// Second call - should call impl again (no caching)
		wrapped.ContextNoError(context.Background(), "a", "b")
		assert.Equal(t, int32(2), atomic.LoadInt32(&impl.contextNoErrorCount))

		// NoParamsOrResults
		wrapped.NoParamsOrResults()
		assert.Equal(t, int32(1), atomic.LoadInt32(&impl.noParamsCallCount))

		wrapped.NoParamsOrResults()
		assert.Equal(t, int32(2), atomic.LoadInt32(&impl.noParamsCallCount))
	})
}
