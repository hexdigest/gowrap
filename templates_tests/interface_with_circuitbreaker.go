package templatestests

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/circuitbreaker template

//go:generate gowrap gen -p github.com/hexdigest/gowrap/templates_tests -i TestInterface -t ../templates/circuitbreaker -o interface_with_circuitbreaker.go

import (
	"context"
	"errors"
	"sync"
	"time"
)

// TestInterfaceWithCircuitBreaker implements TestInterface instrumented with circuit breaker
type TestInterfaceWithCircuitBreaker struct {
	TestInterface

	_lock                 sync.RWMutex
	_maxConsecutiveErrors int
	_consecutiveErrors    int
	_openInterval         time.Duration
	_closesAt             *time.Time
}

// NewTestInterfaceWithCircuitBreaker breakes a circuit after consecutiveErrors of errors and opens the circuit again after openInterval of time.
// If after openInterval first method call results in error we close open again.
func NewTestInterfaceWithCircuitBreaker(base TestInterface, consecutiveErrors int, openInterval time.Duration) *TestInterfaceWithCircuitBreaker {
	return &TestInterfaceWithCircuitBreaker{
		TestInterface:         base,
		_maxConsecutiveErrors: consecutiveErrors,
		_openInterval:         openInterval,
	}
}

// Channels implements TestInterface

// F implements TestInterface
func (_d *TestInterfaceWithCircuitBreaker) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	_d._lock.RLock()

	if _d._closesAt != nil && _d._closesAt.After(time.Now()) {
		_d._lock.RUnlock()
		err = errors.New("TestInterfaceWithCircuitBreaker: circuit is open")
		return
	}
	_d._lock.RUnlock()

	result1, result2, err = _d.TestInterface.F(ctx, a1, a2...)
	_d._lock.Lock()
	defer _d._lock.Unlock()

	if err == nil {
		_d._consecutiveErrors = 0
		_d._closesAt = nil
		return
	}

	_d._consecutiveErrors++

	if _d._consecutiveErrors >= _d._maxConsecutiveErrors {
		closesAt := time.Now().Add(_d._openInterval)
		_d._closesAt = &closesAt
	}

	return
}

// NoError implements TestInterface

// NoParamsOrResults implements TestInterface
