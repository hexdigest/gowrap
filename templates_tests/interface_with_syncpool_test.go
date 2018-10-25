package templatestests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTestInterfacePool(t *testing.T) {
	assert.Panics(t, func() {
		NewTestInterfacePool()
	})
}

func TestTestInterfacePool_F(t *testing.T) {
	start := time.Now()
	delay := 10 * time.Millisecond

	impl1 := &testImpl{delay: delay}
	impl2 := &testImpl{delay: delay}

	wrapped := NewTestInterfacePool(impl1, impl2)

	doneCh1 := make(chan struct{})
	doneCh2 := make(chan struct{})
	doneCh3 := make(chan struct{})
	go func() {
		wrapped.F(context.Background(), "a1", "a2")
		close(doneCh1)
	}()
	go func() {
		wrapped.F(context.Background(), "a1", "a2")
		close(doneCh2)
	}()
	go func() {
		wrapped.F(context.Background(), "a1", "a2")
		close(doneCh3)
	}()

	<-doneCh1
	<-doneCh2
	<-doneCh3

	assert.True(t, time.Now().Sub(start) > delay*2)
	assert.EqualValues(t, 3, impl1.callCounter+impl2.callCounter)
}
