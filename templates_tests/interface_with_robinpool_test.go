package templatestests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestInterfaceRoundRobinPool_F(t *testing.T) {
	impl1 := &testImpl{r1: "11", r2: "12"}
	impl2 := &testImpl{r1: "21", r2: "22"}

	wrapped, err := NewTestInterfaceRoundRobinPool(impl1, impl2)
	require.NoError(t, err)

	for i := 0; i < 8; i++ {
		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")
		assert.NoError(t, err)
		if i%2 == 0 {
			assert.Equal(t, "21", r1)
			assert.Equal(t, "22", r2)
		} else {
			assert.Equal(t, "11", r1)
			assert.Equal(t, "12", r2)
		}
	}

	assert.EqualValues(t, 4, impl1.callCounter)
	assert.EqualValues(t, 4, impl2.callCounter)
}

func TestNewTestInterfaceRoundRobinPool(t *testing.T) {
	wrapped, err := NewTestInterfaceRoundRobinPool()
	require.Error(t, err)
	require.Nil(t, wrapped)
}

func TestMustNewTestInterfaceRoundRobinPool(t *testing.T) {
	assert.Panics(t, func() {
		MustNewTestInterfaceRoundRobinPool()
	})

	assert.NotPanics(t, func() {
		MustNewTestInterfaceRoundRobinPool(&testImpl{})
	})
}
