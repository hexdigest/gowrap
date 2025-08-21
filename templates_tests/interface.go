package templatestests

import "context"

// TestInterface is used to test templates
type TestInterface interface {
	F(ctx context.Context, a1 string, a2 ...string) (result1, result2 string, err error)
	ContextNoError(ctx context.Context, a1 string, a2 string)
	NoError(string) string
	NoParamsOrResults()
	Channels(chA chan bool, chB chan<- bool, chanC <-chan bool)
}

// AnotherTestInterface is used to test templates where TestInterface was already used
type AnotherTestInterface interface {
	F(ctx context.Context, a1 string, a2 ...string) (result1, result2 string, err error)
	NoError(string) string
	NoParamsOrResults()
	Channels(chA chan bool, chB chan<- bool, chanC <-chan bool)
}

// GenericsTestInterface is a generic interface used to test templates with generics
type GenericsTestInterface[T any, U TestInterface] interface {
	F(ctx context.Context, a1 T, a2 ...U) (result1 T, result2 string, err error)
	G(a1 U, a2 string) (result1 T, result2 U, err error)
	ContextNoError(ctx context.Context, a1 string, a2 string)
	NoError(string) string
	NoParamsOrResults()
	Channels(chA chan bool, chB chan<- bool, chanC <-chan bool)
}
