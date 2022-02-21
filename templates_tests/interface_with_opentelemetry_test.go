package templatestests

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func TestTestInterfaceWithOpenTelemetryTracing_F(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2"}
		wrapped := NewTestInterfaceWithOpentelemetry(impl, "test")

		mc := minimock.NewController(t)
		defer mc.Finish()

		span := NewOpentelemetrySpanMock(mc)
		span.EndMock.Expect().Return()

		tr := NewOpentelemetryTracerMock(mc)
		tr.StartMock.Expect(context.Background(), "TestInterface.F").Return(context.Background(), span)

		tp := NewOpentelemetryTracerProviderMock(mc)
		tp.TracerMock.Expect("test").Return(tr)

		otel.SetTracerProvider(tp)

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")
		assert.NoError(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)
	})

	t.Run("error", func(t *testing.T) {
		err := errors.New("unexpected error")

		impl := &testImpl{r1: "1", r2: "2", err: err}
		wrapped := NewTestInterfaceWithOpentelemetry(impl, "test")

		mc := minimock.NewController(t)
		defer mc.Finish()

		span := NewOpentelemetrySpanMock(mc)
		span.RecordErrorMock.Expect(err)
		span.EndMock.Expect().Return()
		span.SetAttributesMock.
			Inspect(func(kv ...attribute.KeyValue) {
				assert.Equal(t, string(kv[0].Key), "event")
				assert.Equal(t, kv[0].Value.AsString(), "error")
				assert.Equal(t, string(kv[1].Key), "message")
				assert.Equal(t, kv[1].Value.AsString(), err.Error())
			}).Return()

		tr := NewOpentelemetryTracerMock(mc)
		tr.StartMock.Expect(context.Background(), "TestInterface.F").Return(context.Background(), span)

		tp := NewOpentelemetryTracerProviderMock(mc)
		tp.TracerMock.Expect("test").Return(tr)

		otel.SetTracerProvider(tp)

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")

		assert.Error(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)

	})
}
