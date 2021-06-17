package templatestests

import (
	"context"
	"errors"
	"go.opencensus.io/trace"
	"testing"

	minimock "github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestTestInterfaceWithOpenCensusTracing_F(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2"}
		wrapped := NewTestInterfaceWithOpenCensus(impl, "test")

		mc := minimock.NewController(t)
		defer mc.Finish()

		span := NewSpanInterfaceMock(mc)
		span.EndMock.Return()

		tracer := NewTracerOpenCensusMock(mc)
		tracer.StartSpanMock.Return(context.Background(), trace.NewSpan(span))

		trace.DefaultTracer = tracer

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")
		assert.NoError(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)
	})

	t.Run("error", func(t *testing.T) {
		err := errors.New("unexpected error")

		impl := &testImpl{r1: "1", r2: "2", err: err}
		wrapped := NewTestInterfaceWithOpenCensus(impl, "test")

		mc := minimock.NewController(t)
		defer mc.Finish()

		span := NewSpanInterfaceMock(mc)
		span.AddAttributesMock.Expect(
			trace.BoolAttribute("error", true),
			trace.StringAttribute("event", "error"),
			trace.StringAttribute("message", err.Error()),
		)
		span.IsRecordingEventsMock.Return(true)
		span.EndMock.Return()

		tracer := NewTracerOpenCensusMock(mc)
		tracer.StartSpanMock.Return(context.Background(), trace.NewSpan(span))

		trace.DefaultTracer = tracer

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")

		assert.Error(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)

	})
}
