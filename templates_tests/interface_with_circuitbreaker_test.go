package templatestests

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type consecutiveErrorsImpl struct {
	NumErrors    int
	NumSuccesses int

	errors    int32
	successes int32
}

var errConsecutive = errors.New("consecutive")

func (c *consecutiveErrorsImpl) F(ctx context.Context, a1 string, a2 ...string) (r1, r2 string, err error) {
	if atomic.AddInt32(&c.errors, 1) > int32(c.NumErrors) {
		if atomic.AddInt32(&c.successes, 1) <= int32(c.NumSuccesses) {
			return "", "", nil
		}

		c.errors = 0
		c.successes = 0
	}

	return "", "", errConsecutive
}

func TestTestInterfaceWithCircuitBreaker_F(t *testing.T) {
	ctx := context.Background()

	t.Run("circuit opens", func(t *testing.T) {
		impl := &consecutiveErrorsImpl{NumErrors: 2, NumSuccesses: 0}
		wrapped := NewTestInterfaceWithCircuitBreaker(impl, 2, time.Second)

		_, _, err := wrapped.F(ctx, "")
		assert.Equal(t, errConsecutive, err)

		_, _, err = wrapped.F(ctx, "")
		assert.Equal(t, errConsecutive, err)

		_, _, err = wrapped.F(ctx, "")
		assert.Equal(t, "TestInterfaceWithCircuitBreaker: circuit is open", err.Error())
	})

	t.Run("circuit closes after open interval", func(t *testing.T) {
		impl := &consecutiveErrorsImpl{NumErrors: 2, NumSuccesses: 10}
		wrapped := NewTestInterfaceWithCircuitBreaker(impl, 2, time.Millisecond)

		_, _, err := wrapped.F(ctx, "")
		assert.Equal(t, errConsecutive, err)

		_, _, err = wrapped.F(ctx, "")
		assert.Equal(t, errConsecutive, err)

		time.Sleep(2 * time.Millisecond)

		_, _, err = wrapped.F(ctx, "")
		assert.NoError(t, err)
	})

	t.Run("circuit opens and closes again", func(t *testing.T) {
		impl := &consecutiveErrorsImpl{NumErrors: 10, NumSuccesses: 0}
		wrapped := NewTestInterfaceWithCircuitBreaker(impl, 2, time.Millisecond)

		_, _, err := wrapped.F(ctx, "")
		assert.Equal(t, errConsecutive, err)

		_, _, err = wrapped.F(ctx, "")
		assert.Equal(t, errConsecutive, err)

		_, _, err = wrapped.F(ctx, "")
		assert.Equal(t, "TestInterfaceWithCircuitBreaker: circuit is open", err.Error())

		time.Sleep(2 * time.Millisecond)

		_, _, err = wrapped.F(ctx, "")
		assert.Equal(t, errConsecutive, err)

		_, _, err = wrapped.F(ctx, "")
		assert.Equal(t, "TestInterfaceWithCircuitBreaker: circuit is open", err.Error())
	})
}
