package templatestests

import "context"

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/syncpool template

//go:generate gowrap gen -p github.com/hexdigest/gowrap/templates_tests -i TestInterface -t ../templates/syncpool -o interface_with_syncpool.go

// TestInterfacePool implements TestInterface that uses pool of TestInterface
type TestInterfacePool struct {
	pool chan TestInterface
}

// NewTestInterfacePool takes several implementations of the TestInterface and returns an instance of the TestInterface
// that uses sync.Pool of given implemetations
func NewTestInterfacePool(impls ...TestInterface) TestInterfacePool {
	if len(impls) == 0 {
		panic("empty pool")
	}

	pool := make(chan TestInterface, len(impls))
	for _, i := range impls {
		pool <- i
	}

	return TestInterfacePool{pool: pool}
}

// F implements TestInterface
func (_d TestInterfacePool) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	_impl := <-_d.pool
	defer func() {
		_d.pool <- _impl
	}()
	return _impl.F(ctx, a1, a2...)
}

// NoError implements TestInterface
func (_d TestInterfacePool) NoError(s1 string) (s2 string) {
	_impl := <-_d.pool
	defer func() {
		_d.pool <- _impl
	}()
	return _impl.NoError(s1)
}

// NoParamsOrResults implements TestInterface
func (_d TestInterfacePool) NoParamsOrResults() {
	_impl := <-_d.pool
	defer func() {
		_d.pool <- _impl
	}()
	_impl.NoParamsOrResults()
	return
}
