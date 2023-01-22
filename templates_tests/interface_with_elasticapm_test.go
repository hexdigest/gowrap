package templatestests

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"go.elastic.co/apm/v2"
)

type ElasticAPM interface {
	StartSpan(ctx context.Context, name, spanType string) (*apm.Span, context.Context)
	EndSpan(span *apm.Span)
	SetLabel(span *apm.Span, key string, value interface{})
	CaptureError(ctx context.Context, err error)
}

func reintroduceElasticAPM(testAPMTracing *TestInterfaceAPMTracing, apm ElasticAPM) {
	testAPMTracing.startSpan = apm.StartSpan
	testAPMTracing.endSpan = apm.EndSpan
	testAPMTracing.setLabel = apm.SetLabel
	testAPMTracing.captureError = apm.CaptureError
}

func TestTestInterfaceWithElasticAPMTracing_F(t *testing.T) {
	t.Run("F no error", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2"}

		wrapped := NewTestInterfaceAPMTracing(impl)

		mc := minimock.NewController(t)
		defer mc.Finish()

		elasticAPM := NewElasticAPMMock(mc)
		ctx := context.Background()
		span, ctxSpan := apm.StartSpan(ctx, "testinterface.F", "testinterface")
		defer span.End()

		elasticAPM.
			StartSpanMock.Expect(ctx, "testinterface.F", "testinterface").Return(span, ctxSpan).
			SetLabelMock.Return().
			EndSpanMock.Expect(span).Return()
		reintroduceElasticAPM(&wrapped, elasticAPM)

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")
		assert.NoError(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)
	})

	t.Run("F error", func(t *testing.T) {
		err := errors.New("unexpected error")

		impl := &testImpl{r1: "1", r2: "2", err: err}
		wrapped := NewTestInterfaceAPMTracing(impl)

		mc := minimock.NewController(t)
		defer mc.Finish()

		elasticAPM := NewElasticAPMMock(mc)
		ctx := context.Background()
		span, ctxSpan := apm.StartSpan(ctx, "testinterface.F", "testinterface")
		defer span.End()

		elasticAPM.
			StartSpanMock.Expect(ctx, "testinterface.F", "testinterface").Return(span, ctxSpan).
			SetLabelMock.Return().
			CaptureErrorMock.Expect(ctxSpan, err).Return().
			EndSpanMock.Expect(span).Return()
		reintroduceElasticAPM(&wrapped, elasticAPM)

		r1, r2, err := wrapped.F(context.Background(), "a1", "a2")

		assert.Error(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)
	})

	t.Run("ContextNoError", func(t *testing.T) {
		err := errors.New("unexpected error")

		impl := &testImpl{r1: "1", r2: "2", err: err}
		wrapped := NewTestInterfaceAPMTracing(impl)

		mc := minimock.NewController(t)
		defer mc.Finish()

		elasticAPM := NewElasticAPMMock(mc)
		ctx := context.Background()
		span, ctxSpan := apm.StartSpan(ctx, "testinterface.ContextNoError", "testinterface")
		defer span.End()

		elasticAPM.
			StartSpanMock.Expect(ctx, "testinterface.ContextNoError", "testinterface").Return(span, ctxSpan).
			SetLabelMock.Return().
			EndSpanMock.Expect(span).Return()
		reintroduceElasticAPM(&wrapped, elasticAPM)

		wrapped.ContextNoError(context.Background(), "a1", "a2")
	})
}
