package loader

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "httpClient" can be found in github.com/hexdigest/gowrap/loader
*/
import (
	http "net/http"
	"sync/atomic"
	"time"

	minimock "github.com/gojuno/minimock/v3"

	testify_assert "github.com/stretchr/testify/assert"
)

//httpClientMock implements github.com/hexdigest/gowrap/loader.httpClient
type httpClientMock struct {
	t minimock.Tester

	DoFunc       func(p *http.Request) (r *http.Response, r1 error)
	DoCounter    uint64
	DoPreCounter uint64
	DoMock       mhttpClientMockDo
}

//newHTTPClientMock returns a mock for github.com/hexdigest/gowrap/loader.httpClient
func newHTTPClientMock(t minimock.Tester) *httpClientMock {
	m := &httpClientMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DoMock = mhttpClientMockDo{mock: m}

	return m
}

type mhttpClientMockDo struct {
	mock             *httpClientMock
	mockExpectations *httpClientMockDoParams
}

//httpClientMockDoParams represents input parameters of the httpClient.Do
type httpClientMockDoParams struct {
	p *http.Request
}

//Expect sets up expected params for the httpClient.Do
func (m *mhttpClientMockDo) Expect(p *http.Request) *mhttpClientMockDo {
	m.mockExpectations = &httpClientMockDoParams{p}
	return m
}

//Return sets up a mock for httpClient.Do to return Return's arguments
func (m *mhttpClientMockDo) Return(r *http.Response, r1 error) *httpClientMock {
	m.mock.DoFunc = func(p *http.Request) (*http.Response, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of httpClient.Do method
func (m *mhttpClientMockDo) Set(f func(p *http.Request) (r *http.Response, r1 error)) *httpClientMock {
	m.mock.DoFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Do implements github.com/hexdigest/gowrap/loader.httpClient interface
func (m *httpClientMock) Do(p *http.Request) (r *http.Response, r1 error) {
	atomic.AddUint64(&m.DoPreCounter, 1)
	defer atomic.AddUint64(&m.DoCounter, 1)

	if m.DoMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.DoMock.mockExpectations, httpClientMockDoParams{p},
			"httpClient.Do got unexpected parameters")

		if m.DoFunc == nil {

			m.t.Fatal("No results are set for the httpClientMock.Do")

			return
		}
	}

	if m.DoFunc == nil {
		m.t.Fatal("Unexpected call to httpClientMock.Do")
		return
	}

	return m.DoFunc(p)
}

//DoMinimockCounter returns a count of httpClientMock.DoFunc invocations
func (m *httpClientMock) DoMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DoCounter)
}

//DoMinimockPreCounter returns the value of httpClientMock.Do invocations
func (m *httpClientMock) DoMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DoPreCounter)
}

func (m *httpClientMock) minimockCheckDoCalled(f func(...interface{})) {
	if m.DoFunc != nil && atomic.LoadUint64(&m.DoCounter) == 0 {
		f("Expected call to httpClientMock.Do")
	}
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *httpClientMock) MinimockFinish() {
	m.minimockCheck(m.t.Fatal)
}

func (m *httpClientMock) minimockCheck(f func(...interface{})) {
	m.minimockCheckDoCalled(f)

}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *httpClientMock) MinimockWait(timeout time.Duration) {
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
func (m *httpClientMock) AllMocksCalled() bool {
	result := true
	result = result && (m.DoFunc == nil || atomic.LoadUint64(&m.DoCounter) > 0)

	return result
}
