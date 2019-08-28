package templatestests

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Tracer" can be found in github.com/opentracing/opentracing-go
*/
import (
	"sync/atomic"
	"time"

	minimock "github.com/gojuno/minimock/v3"
	opentracing "github.com/opentracing/opentracing-go"

	testify_assert "github.com/stretchr/testify/assert"
)

//TracerMock implements github.com/opentracing/opentracing-go.Tracer
type TracerMock struct {
	t minimock.Tester

	ExtractFunc       func(p interface{}, p1 interface{}) (r opentracing.SpanContext, r1 error)
	ExtractCounter    uint64
	ExtractPreCounter uint64
	ExtractMock       mTracerMockExtract

	InjectFunc       func(p opentracing.SpanContext, p1 interface{}, p2 interface{}) (r error)
	InjectCounter    uint64
	InjectPreCounter uint64
	InjectMock       mTracerMockInject

	StartSpanFunc       func(p string, p1 ...opentracing.StartSpanOption) (r opentracing.Span)
	StartSpanCounter    uint64
	StartSpanPreCounter uint64
	StartSpanMock       mTracerMockStartSpan
}

//NewTracerMock returns a mock for github.com/opentracing/opentracing-go.Tracer
func NewTracerMock(t minimock.Tester) *TracerMock {
	m := &TracerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExtractMock = mTracerMockExtract{mock: m}
	m.InjectMock = mTracerMockInject{mock: m}
	m.StartSpanMock = mTracerMockStartSpan{mock: m}

	return m
}

type mTracerMockExtract struct {
	mock             *TracerMock
	mockExpectations *TracerMockExtractParams
}

//TracerMockExtractParams represents input parameters of the Tracer.Extract
type TracerMockExtractParams struct {
	p  interface{}
	p1 interface{}
}

//Expect sets up expected params for the Tracer.Extract
func (m *mTracerMockExtract) Expect(p interface{}, p1 interface{}) *mTracerMockExtract {
	m.mockExpectations = &TracerMockExtractParams{p, p1}
	return m
}

//Return sets up a mock for Tracer.Extract to return Return's arguments
func (m *mTracerMockExtract) Return(r opentracing.SpanContext, r1 error) *TracerMock {
	m.mock.ExtractFunc = func(p interface{}, p1 interface{}) (opentracing.SpanContext, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Tracer.Extract method
func (m *mTracerMockExtract) Set(f func(p interface{}, p1 interface{}) (r opentracing.SpanContext, r1 error)) *TracerMock {
	m.mock.ExtractFunc = f
	return m.mock
}

//Extract implements github.com/opentracing/opentracing-go.Tracer interface
func (m *TracerMock) Extract(p interface{}, p1 interface{}) (r opentracing.SpanContext, r1 error) {
	atomic.AddUint64(&m.ExtractPreCounter, 1)
	defer atomic.AddUint64(&m.ExtractCounter, 1)

	if m.ExtractMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExtractMock.mockExpectations, TracerMockExtractParams{p, p1},
			"Tracer.Extract got unexpected parameters")

		if m.ExtractFunc == nil {

			m.t.Fatal("No results are set for the TracerMock.Extract")

			return
		}
	}

	if m.ExtractFunc == nil {
		m.t.Fatal("Unexpected call to TracerMock.Extract")
		return
	}

	return m.ExtractFunc(p, p1)
}

//ExtractMinimockCounter returns a count of TracerMock.ExtractFunc invocations
func (m *TracerMock) ExtractMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExtractCounter)
}

//ExtractMinimockPreCounter returns the value of TracerMock.Extract invocations
func (m *TracerMock) ExtractMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExtractPreCounter)
}

type mTracerMockInject struct {
	mock             *TracerMock
	mockExpectations *TracerMockInjectParams
}

//TracerMockInjectParams represents input parameters of the Tracer.Inject
type TracerMockInjectParams struct {
	p  opentracing.SpanContext
	p1 interface{}
	p2 interface{}
}

//Expect sets up expected params for the Tracer.Inject
func (m *mTracerMockInject) Expect(p opentracing.SpanContext, p1 interface{}, p2 interface{}) *mTracerMockInject {
	m.mockExpectations = &TracerMockInjectParams{p, p1, p2}
	return m
}

//Return sets up a mock for Tracer.Inject to return Return's arguments
func (m *mTracerMockInject) Return(r error) *TracerMock {
	m.mock.InjectFunc = func(p opentracing.SpanContext, p1 interface{}, p2 interface{}) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Tracer.Inject method
func (m *mTracerMockInject) Set(f func(p opentracing.SpanContext, p1 interface{}, p2 interface{}) (r error)) *TracerMock {
	m.mock.InjectFunc = f
	return m.mock
}

//Inject implements github.com/opentracing/opentracing-go.Tracer interface
func (m *TracerMock) Inject(p opentracing.SpanContext, p1 interface{}, p2 interface{}) (r error) {
	atomic.AddUint64(&m.InjectPreCounter, 1)
	defer atomic.AddUint64(&m.InjectCounter, 1)

	if m.InjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.InjectMock.mockExpectations, TracerMockInjectParams{p, p1, p2},
			"Tracer.Inject got unexpected parameters")

		if m.InjectFunc == nil {

			m.t.Fatal("No results are set for the TracerMock.Inject")

			return
		}
	}

	if m.InjectFunc == nil {
		m.t.Fatal("Unexpected call to TracerMock.Inject")
		return
	}

	return m.InjectFunc(p, p1, p2)
}

//InjectMinimockCounter returns a count of TracerMock.InjectFunc invocations
func (m *TracerMock) InjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InjectCounter)
}

//InjectMinimockPreCounter returns the value of TracerMock.Inject invocations
func (m *TracerMock) InjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InjectPreCounter)
}

type mTracerMockStartSpan struct {
	mock             *TracerMock
	mockExpectations *TracerMockStartSpanParams
}

//TracerMockStartSpanParams represents input parameters of the Tracer.StartSpan
type TracerMockStartSpanParams struct {
	p  string
	p1 []opentracing.StartSpanOption
}

//Expect sets up expected params for the Tracer.StartSpan
func (m *mTracerMockStartSpan) Expect(p string, p1 ...opentracing.StartSpanOption) *mTracerMockStartSpan {
	m.mockExpectations = &TracerMockStartSpanParams{p, p1}
	return m
}

//Return sets up a mock for Tracer.StartSpan to return Return's arguments
func (m *mTracerMockStartSpan) Return(r opentracing.Span) *TracerMock {
	m.mock.StartSpanFunc = func(p string, p1 ...opentracing.StartSpanOption) opentracing.Span {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Tracer.StartSpan method
func (m *mTracerMockStartSpan) Set(f func(p string, p1 ...opentracing.StartSpanOption) (r opentracing.Span)) *TracerMock {
	m.mock.StartSpanFunc = f
	return m.mock
}

//StartSpan implements github.com/opentracing/opentracing-go.Tracer interface
func (m *TracerMock) StartSpan(p string, p1 ...opentracing.StartSpanOption) (r opentracing.Span) {
	atomic.AddUint64(&m.StartSpanPreCounter, 1)
	defer atomic.AddUint64(&m.StartSpanCounter, 1)

	if m.StartSpanMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.StartSpanMock.mockExpectations, TracerMockStartSpanParams{p, p1},
			"Tracer.StartSpan got unexpected parameters")

		if m.StartSpanFunc == nil {

			m.t.Fatal("No results are set for the TracerMock.StartSpan")

			return
		}
	}

	if m.StartSpanFunc == nil {
		m.t.Fatal("Unexpected call to TracerMock.StartSpan")
		return
	}

	return m.StartSpanFunc(p, p1...)
}

//StartSpanMinimockCounter returns a count of TracerMock.StartSpanFunc invocations
func (m *TracerMock) StartSpanMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StartSpanCounter)
}

//StartSpanMinimockPreCounter returns the value of TracerMock.StartSpan invocations
func (m *TracerMock) StartSpanMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StartSpanPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TracerMock) ValidateCallCounters() {

	if m.ExtractFunc != nil && atomic.LoadUint64(&m.ExtractCounter) == 0 {
		m.t.Fatal("Expected call to TracerMock.Extract")
	}

	if m.InjectFunc != nil && atomic.LoadUint64(&m.InjectCounter) == 0 {
		m.t.Fatal("Expected call to TracerMock.Inject")
	}

	if m.StartSpanFunc != nil && atomic.LoadUint64(&m.StartSpanCounter) == 0 {
		m.t.Fatal("Expected call to TracerMock.StartSpan")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TracerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TracerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TracerMock) MinimockFinish() {

	if m.ExtractFunc != nil && atomic.LoadUint64(&m.ExtractCounter) == 0 {
		m.t.Fatal("Expected call to TracerMock.Extract")
	}

	if m.InjectFunc != nil && atomic.LoadUint64(&m.InjectCounter) == 0 {
		m.t.Fatal("Expected call to TracerMock.Inject")
	}

	if m.StartSpanFunc != nil && atomic.LoadUint64(&m.StartSpanCounter) == 0 {
		m.t.Fatal("Expected call to TracerMock.StartSpan")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TracerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *TracerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.ExtractFunc == nil || atomic.LoadUint64(&m.ExtractCounter) > 0)
		ok = ok && (m.InjectFunc == nil || atomic.LoadUint64(&m.InjectCounter) > 0)
		ok = ok && (m.StartSpanFunc == nil || atomic.LoadUint64(&m.StartSpanCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ExtractFunc != nil && atomic.LoadUint64(&m.ExtractCounter) == 0 {
				m.t.Error("Expected call to TracerMock.Extract")
			}

			if m.InjectFunc != nil && atomic.LoadUint64(&m.InjectCounter) == 0 {
				m.t.Error("Expected call to TracerMock.Inject")
			}

			if m.StartSpanFunc != nil && atomic.LoadUint64(&m.StartSpanCounter) == 0 {
				m.t.Error("Expected call to TracerMock.StartSpan")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *TracerMock) AllMocksCalled() bool {

	if m.ExtractFunc != nil && atomic.LoadUint64(&m.ExtractCounter) == 0 {
		return false
	}

	if m.InjectFunc != nil && atomic.LoadUint64(&m.InjectCounter) == 0 {
		return false
	}

	if m.StartSpanFunc != nil && atomic.LoadUint64(&m.StartSpanCounter) == 0 {
		return false
	}

	return true
}
