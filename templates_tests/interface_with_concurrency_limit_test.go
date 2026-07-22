package templatestests

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTestInterfaceWithConcurrencyLimit_F(t *testing.T) {
	impl := &testImpl{r1: "1", r2: "2", delay: 100 * time.Millisecond}

	wrapped := NewTestInterfaceWithConcurrencyLimit(impl, 3)

	for i := 0; i < 10; i++ {
		go func() {
			r1, r2, err := wrapped.F(context.Background(), "a1")
			assert.NoError(t, err)
			assert.Equal(t, "1", r1)
			assert.Equal(t, "2", r2)

		}()
	}

	<-time.After(10 * time.Millisecond)

	counter := atomic.LoadUint64(&impl.callCounter)
	assert.EqualValues(t, 3, counter) // the first burst

	<-time.After(100 * time.Millisecond)

	counter = atomic.LoadUint64(&impl.callCounter)
	assert.EqualValues(t, 6, counter) // the second burst

	<-time.After(100 * time.Millisecond)

	counter = atomic.LoadUint64(&impl.callCounter)
	assert.EqualValues(t, 9, counter) // the third burst

	<-time.After(100 * time.Millisecond)

	counter = atomic.LoadUint64(&impl.callCounter)
	assert.EqualValues(t, 10, counter) // the 10th call
}
