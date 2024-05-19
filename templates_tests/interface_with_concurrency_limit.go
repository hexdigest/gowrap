package templatestests

import "context"

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/concurrencylimit template

//go:generate gowrap gen -p github.com/hexdigest/gowrap/templates_tests -i TestInterface -t ../templates/concurrencylimit -o interface_with_concurrency_limit.go

// TestInterfaceWithConcurrencyLimit implements TestInterface
type TestInterfaceWithConcurrencyLimit struct {
	_base  TestInterface
	_burst chan int
}

// NewTestInterfaceWithConcurrencyLimit instruments an implementation of the TestInterface with concurrency limiting
func NewTestInterfaceWithConcurrencyLimit(base TestInterface, concurrentCalls int) *TestInterfaceWithConcurrencyLimit {
	d := &TestInterfaceWithConcurrencyLimit{
		_base:  base,
		_burst: make(chan int, concurrentCalls),
	}

	return d
}

// Channels implements TestInterface
func (_d *TestInterfaceWithConcurrencyLimit) Channels(chA chan bool, chB chan<- bool, chanC <-chan bool) {
	_d._burst <- 1
	defer func() {
		<-_d._burst
	}()

	_d._base.Channels(chA, chB, chanC)
	return
}

// F implements TestInterface
func (_d *TestInterfaceWithConcurrencyLimit) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	case _d._burst <- 1:
		defer func() {
			<-_d._burst
		}()
	}

	return _d._base.F(ctx, a1, a2...)
}

// NoError implements TestInterface
func (_d *TestInterfaceWithConcurrencyLimit) NoError(s1 string) (s2 string) {
	_d._burst <- 1
	defer func() {
		<-_d._burst
	}()

	return _d._base.NoError(s1)
}

// NoParamsOrResults implements TestInterface
func (_d *TestInterfaceWithConcurrencyLimit) NoParamsOrResults() {
	_d._burst <- 1
	defer func() {
		<-_d._burst
	}()

	_d._base.NoParamsOrResults()
	return
}
