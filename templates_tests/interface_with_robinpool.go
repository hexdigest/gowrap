package templatestests

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/robinpool template

//go:generate gowrap gen -d . -i TestInterface -t ../templates/robinpool -o interface_with_robinpool.go

import (
	"context"
	"errors"
	"sync/atomic"
)

// TestInterfaceRoundRobinPool implements TestInterface that uses pool of TestInterface
type TestInterfaceRoundRobinPool struct {
	pool     []TestInterface
	poolSize uint32
	counter  uint32
}

// NewTestInterfaceRoundRobinPool takes several implementations of the TestInterface and returns an instance of the TestInterface
// that picks one of the given implementations using Round-robin algorithm and delegates method call to it
func NewTestInterfaceRoundRobinPool(pool ...TestInterface) (*TestInterfaceRoundRobinPool, error) {
	if len(pool) == 0 {
		return nil, errors.New("empty pool")
	}

	return &TestInterfaceRoundRobinPool{pool: pool, poolSize: uint32(len(pool))}, nil
}

// MustNewTestInterfaceRoundRobinPool takes several implementations of the TestInterface and returns an instance of the TestInterface
// that picks one of the given implementations using Round-robin algorithm and delegates method call to it.
func MustNewTestInterfaceRoundRobinPool(pool ...TestInterface) *TestInterfaceRoundRobinPool {
	if len(pool) == 0 {
		panic("empty pool")
	}

	return &TestInterfaceRoundRobinPool{pool: pool, poolSize: uint32(len(pool))}
}

// F implements TestInterface
func (_d *TestInterfaceRoundRobinPool) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	_counter := atomic.AddUint32(&_d.counter, 1)
	return _d.pool[_counter%_d.poolSize].F(ctx, a1, a2...)
}
