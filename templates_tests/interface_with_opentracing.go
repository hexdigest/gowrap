package templatestests

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/opentracing template

//go:generate gowrap gen -p github.com/hexdigest/gowrap/templates_tests -i TestInterface -t ../templates/opentracing -o interface_with_opentracing.go

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	_ext "github.com/opentracing/opentracing-go/ext"
	_log "github.com/opentracing/opentracing-go/log"
)

// TestInterfaceWithTracing implements TestInterface interface instrumented with opentracing spans
type TestInterfaceWithTracing struct {
	TestInterface
	_instance      string
	_spanDecorator func(span opentracing.Span, params, results map[string]interface{})
}

// NewTestInterfaceWithTracing returns TestInterfaceWithTracing
func NewTestInterfaceWithTracing(base TestInterface, instance string, spanDecorator ...func(span opentracing.Span, params, results map[string]interface{})) TestInterfaceWithTracing {
	d := TestInterfaceWithTracing{
		TestInterface: base,
		_instance:     instance,
	}

	if len(spanDecorator) > 0 && spanDecorator[0] != nil {
		d._spanDecorator = spanDecorator[0]
	}

	return d
}

// F implements TestInterface
func (_d TestInterfaceWithTracing) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	_span, ctx := opentracing.StartSpanFromContext(ctx, _d._instance+".TestInterface.F")
	defer func() {
		if _d._spanDecorator != nil {
			_d._spanDecorator(_span, map[string]interface{}{
				"ctx": ctx,
				"a1":  a1,
				"a2":  a2}, map[string]interface{}{
				"result1": result1,
				"result2": result2,
				"err":     err})
		} else if err != nil {
			_ext.Error.Set(_span, true)
			_span.LogFields(
				_log.String("event", "error"),
				_log.String("message", err.Error()),
			)
		}

		_span.Finish()
	}()
	return _d.TestInterface.F(ctx, a1, a2...)
}
