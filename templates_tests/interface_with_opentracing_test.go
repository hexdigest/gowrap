package templatestests

import (
	"context"
	"errors"
	"testing"

	minimock "github.com/gojuno/minimock/v3"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestTestInterfaceWithTracing_F(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2"}
		wrapped := NewTestInterfaceWithTracing(impl, "test")

		mc := minimock.NewController(t)
		defer mc.Finish()

		span := NewSpanMock(mc)
		span.FinishMock.Return()

		tracer := NewTracerMock(mc)
		tracer.StartSpanMock.Return(span)

		opentracing.InitGlobalTracer(tracer)

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")
		assert.NoError(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)
	})

	t.Run("error", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2", err: errors.New("unexpected error")}
		wrapped := NewTestInterfaceWithTracing(impl, "test")

		mc := minimock.NewController(t)
		defer mc.Finish()

		span := NewSpanMock(mc)
		span.FinishMock.Return()
		span.SetTagMock.Return(span)
		span.LogFieldsMock.Return()

		tracer := NewTracerMock(mc)
		tracer.StartSpanMock.Return(span)

		opentracing.InitGlobalTracer(tracer)

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")

		assert.Error(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)

	})
}
