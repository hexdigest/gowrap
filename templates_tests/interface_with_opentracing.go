package templatestests

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/opentracing template

//go:generate gowrap gen -d . -i TestInterface -t ../templates/opentracing -o interface_with_opentracing.go

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	_ext "github.com/opentracing/opentracing-go/ext"
	_log "github.com/opentracing/opentracing-go/log"
)

// TestInterfaceWithTracing implements TestInterface interface instrumented with opentracing spans
type TestInterfaceWithTracing struct {
	TestInterface
	_instance string
}

// NewTestInterfaceWithTracing returns TestInterfaceWithTracing
func NewTestInterfaceWithTracing(base TestInterface, instance string) TestInterfaceWithTracing {
	return TestInterfaceWithTracing{
		TestInterface: base,
		_instance:     instance,
	}
}

// F implements TestInterface
func (_d TestInterfaceWithTracing) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	_span, ctx := opentracing.StartSpanFromContext(ctx, _d._instance+".TestInterface.F")
	defer func() {
		if err != nil {
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
