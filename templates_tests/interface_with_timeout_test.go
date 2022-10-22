package templatestests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type timeoutsImpl struct {
	hasTimeout bool
	t          *testing.T
}

func (c *timeoutsImpl) F(ctx context.Context, a1 string, a2 ...string) (r1, r2 string, err error) {
	_, ok := ctx.Deadline()
	assert.True(c.t, ok == c.hasTimeout)

	return "", "", nil
}

func (c *timeoutsImpl) ContextNoError(ctx context.Context, a1 string, a2 string) {
}

func (c *timeoutsImpl) NoError(string) string {
	return ""
}

func (c *timeoutsImpl) NoParamsOrResults() {
}

func (c *timeoutsImpl) Channels(chA chan bool, chB chan<- bool, chanC <-chan bool) {
}

func TestTestInterfaceWithTimeout_F(t *testing.T) {
	ctx := context.Background()

	t.Run("timeout is not set", func(t *testing.T) {
		impl := &timeoutsImpl{hasTimeout: false, t: t}
		wrapped := NewTestInterfaceWithTimeout(impl, TestInterfaceWithTimeoutConfig{})

		_, _, err := wrapped.F(ctx, "")
		assert.NoError(t, err)
	})

	t.Run("timeout is set", func(t *testing.T) {
		impl := &timeoutsImpl{hasTimeout: true, t: t}
		wrapped := NewTestInterfaceWithTimeout(impl, TestInterfaceWithTimeoutConfig{
			FTimeout: time.Second,
		})

		_, _, err := wrapped.F(ctx, "")
		assert.NoError(t, err)
	})
}
