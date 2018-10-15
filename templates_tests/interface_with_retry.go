package templatestests

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/retry template

//go:generate gowrap gen -d . -i TestInterface -t ../templates/retry -o interface_with_retry.go

import (
	"context"
	"time"
)

// TestInterfaceWithRetry implements TestInterface interface instrumented with retries
type TestInterfaceWithRetry struct {
	TestInterface
	_retryCount    int
	_retryInterval time.Duration
}

// NewTestInterfaceWithRetry returns TestInterfaceWithRetry
func NewTestInterfaceWithRetry(base TestInterface, retryCount int, retryInterval time.Duration) TestInterfaceWithRetry {
	return TestInterfaceWithRetry{
		TestInterface:  base,
		_retryCount:    retryCount + 1,
		_retryInterval: retryInterval,
	}
}

// F implements TestInterface
func (_d TestInterfaceWithRetry) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	for _i := 0; _i < _d._retryCount; _i++ {
		result1, result2, err = _d.F(ctx, a1, a2...)
		if err == nil {
			break
		}
		if _d._retryCount > 1 {
			select {
			case <-ctx.Done():
				return result1, result2, err
			case <-time.After(_d._retryInterval):
			}
		}
	}
	return result1, result2, err
}
