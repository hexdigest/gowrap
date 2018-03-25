package gowrap

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "remoteTemplateLoader" can be found in github.com/hexdigest/gowrap
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//remoteTemplateLoaderMock implements github.com/hexdigest/gowrap.remoteTemplateLoader
type remoteTemplateLoaderMock struct {
	t minimock.Tester

	ListFunc       func() (r []string, r1 error)
	ListCounter    uint64
	ListPreCounter uint64
	ListMock       mremoteTemplateLoaderMockList

	LoadFunc       func(p string) (r []byte, r1 string, r2 error)
	LoadCounter    uint64
	LoadPreCounter uint64
	LoadMock       mremoteTemplateLoaderMockLoad
}

//newRemoteTemplateLoaderMock returns a mock for github.com/hexdigest/gowrap.remoteTemplateLoader
func newRemoteTemplateLoaderMock(t minimock.Tester) *remoteTemplateLoaderMock {
	m := &remoteTemplateLoaderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ListMock = mremoteTemplateLoaderMockList{mock: m}
	m.LoadMock = mremoteTemplateLoaderMockLoad{mock: m}

	return m
}

type mremoteTemplateLoaderMockList struct {
	mock *remoteTemplateLoaderMock
}

//Return sets up a mock for remoteTemplateLoader.List to return Return's arguments
func (m *mremoteTemplateLoaderMockList) Return(r []string, r1 error) *remoteTemplateLoaderMock {
	m.mock.ListFunc = func() ([]string, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of remoteTemplateLoader.List method
func (m *mremoteTemplateLoaderMockList) Set(f func() (r []string, r1 error)) *remoteTemplateLoaderMock {
	m.mock.ListFunc = f

	return m.mock
}

//List implements github.com/hexdigest/gowrap.remoteTemplateLoader interface
func (m *remoteTemplateLoaderMock) List() (r []string, r1 error) {
	atomic.AddUint64(&m.ListPreCounter, 1)
	defer atomic.AddUint64(&m.ListCounter, 1)

	if m.ListFunc == nil {
		m.t.Fatal("Unexpected call to remoteTemplateLoaderMock.List")
		return
	}

	return m.ListFunc()
}

//ListMinimockCounter returns a count of remoteTemplateLoaderMock.ListFunc invocations
func (m *remoteTemplateLoaderMock) ListMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ListCounter)
}

//ListMinimockPreCounter returns the value of remoteTemplateLoaderMock.List invocations
func (m *remoteTemplateLoaderMock) ListMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ListPreCounter)
}

func (m *remoteTemplateLoaderMock) minimockCheckListCalled(f func(...interface{})) {
	if m.ListFunc != nil && atomic.LoadUint64(&m.ListCounter) == 0 {
		f("Expected call to remoteTemplateLoaderMock.List")
	}
}

type mremoteTemplateLoaderMockLoad struct {
	mock             *remoteTemplateLoaderMock
	mockExpectations *remoteTemplateLoaderMockLoadParams
}

//remoteTemplateLoaderMockLoadParams represents input parameters of the remoteTemplateLoader.Load
type remoteTemplateLoaderMockLoadParams struct {
	p string
}

//Expect sets up expected params for the remoteTemplateLoader.Load
func (m *mremoteTemplateLoaderMockLoad) Expect(p string) *mremoteTemplateLoaderMockLoad {
	m.mockExpectations = &remoteTemplateLoaderMockLoadParams{p}
	return m
}

//Return sets up a mock for remoteTemplateLoader.Load to return Return's arguments
func (m *mremoteTemplateLoaderMockLoad) Return(r []byte, r1 string, r2 error) *remoteTemplateLoaderMock {
	m.mock.LoadFunc = func(p string) ([]byte, string, error) {
		return r, r1, r2
	}
	return m.mock
}

//Set uses given function f as a mock of remoteTemplateLoader.Load method
func (m *mremoteTemplateLoaderMockLoad) Set(f func(p string) (r []byte, r1 string, r2 error)) *remoteTemplateLoaderMock {
	m.mock.LoadFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Load implements github.com/hexdigest/gowrap.remoteTemplateLoader interface
func (m *remoteTemplateLoaderMock) Load(p string) (r []byte, r1 string, r2 error) {
	atomic.AddUint64(&m.LoadPreCounter, 1)
	defer atomic.AddUint64(&m.LoadCounter, 1)

	if m.LoadMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LoadMock.mockExpectations, remoteTemplateLoaderMockLoadParams{p},
			"remoteTemplateLoader.Load got unexpected parameters")

		if m.LoadFunc == nil {

			m.t.Fatal("No results are set for the remoteTemplateLoaderMock.Load")

			return
		}
	}

	if m.LoadFunc == nil {
		m.t.Fatal("Unexpected call to remoteTemplateLoaderMock.Load")
		return
	}

	return m.LoadFunc(p)
}

//LoadMinimockCounter returns a count of remoteTemplateLoaderMock.LoadFunc invocations
func (m *remoteTemplateLoaderMock) LoadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LoadCounter)
}

//LoadMinimockPreCounter returns the value of remoteTemplateLoaderMock.Load invocations
func (m *remoteTemplateLoaderMock) LoadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LoadPreCounter)
}

func (m *remoteTemplateLoaderMock) minimockCheckLoadCalled(f func(...interface{})) {
	if m.LoadFunc != nil && atomic.LoadUint64(&m.LoadCounter) == 0 {
		f("Expected call to remoteTemplateLoaderMock.Load")
	}
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *remoteTemplateLoaderMock) MinimockFinish() {
	m.minimockCheck(m.t.Fatal)
}

func (m *remoteTemplateLoaderMock) minimockCheck(f func(...interface{})) {
	m.minimockCheckListCalled(f)
	m.minimockCheckLoadCalled(f)

}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *remoteTemplateLoaderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		if m.AllMocksCalled() {
			return
		}

		select {
		case <-timeoutCh:
			m.minimockCheck(m.t.Error)
			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *remoteTemplateLoaderMock) AllMocksCalled() bool {
	result := true
	result = result && (m.ListFunc == nil || atomic.LoadUint64(&m.ListCounter) > 0)
	result = result && (m.LoadFunc == nil || atomic.LoadUint64(&m.LoadCounter) > 0)

	return result
}
