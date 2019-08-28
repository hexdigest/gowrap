package loader

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ReadCloser" can be found in io
*/
import (
	"sync/atomic"
	"time"

	minimock "github.com/gojuno/minimock/v3"
	testify_assert "github.com/stretchr/testify/assert"
)

//ReadCloserMock implements io.ReadCloser
type ReadCloserMock struct {
	t minimock.Tester

	CloseFunc       func() (r error)
	CloseCounter    uint64
	ClosePreCounter uint64
	CloseMock       mReadCloserMockClose

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mReadCloserMockRead
}

//NewReadCloserMock returns a mock for io.ReadCloser
func NewReadCloserMock(t minimock.Tester) *ReadCloserMock {
	m := &ReadCloserMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloseMock = mReadCloserMockClose{mock: m}
	m.ReadMock = mReadCloserMockRead{mock: m}

	return m
}

type mReadCloserMockClose struct {
	mock *ReadCloserMock
}

//Return sets up a mock for ReadCloser.Close to return Return's arguments
func (m *mReadCloserMockClose) Return(r error) *ReadCloserMock {
	m.mock.CloseFunc = func() error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ReadCloser.Close method
func (m *mReadCloserMockClose) Set(f func() (r error)) *ReadCloserMock {
	m.mock.CloseFunc = f

	return m.mock
}

//Close implements io.ReadCloser interface
func (m *ReadCloserMock) Close() (r error) {
	atomic.AddUint64(&m.ClosePreCounter, 1)
	defer atomic.AddUint64(&m.CloseCounter, 1)

	if m.CloseFunc == nil {
		m.t.Fatal("Unexpected call to ReadCloserMock.Close")
		return
	}

	return m.CloseFunc()
}

//CloseMinimockCounter returns a count of ReadCloserMock.CloseFunc invocations
func (m *ReadCloserMock) CloseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloseCounter)
}

//CloseMinimockPreCounter returns the value of ReadCloserMock.Close invocations
func (m *ReadCloserMock) CloseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClosePreCounter)
}

func (m *ReadCloserMock) minimockCheckCloseCalled(f func(...interface{})) {
	if m.CloseFunc != nil && atomic.LoadUint64(&m.CloseCounter) == 0 {
		f("Expected call to ReadCloserMock.Close")
	}
}

type mReadCloserMockRead struct {
	mock             *ReadCloserMock
	mockExpectations *ReadCloserMockReadParams
}

//ReadCloserMockReadParams represents input parameters of the ReadCloser.Read
type ReadCloserMockReadParams struct {
	p []byte
}

//Expect sets up expected params for the ReadCloser.Read
func (m *mReadCloserMockRead) Expect(p []byte) *mReadCloserMockRead {
	m.mockExpectations = &ReadCloserMockReadParams{p}
	return m
}

//Return sets up a mock for ReadCloser.Read to return Return's arguments
func (m *mReadCloserMockRead) Return(r int, r1 error) *ReadCloserMock {
	m.mock.ReadFunc = func(p []byte) (int, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ReadCloser.Read method
func (m *mReadCloserMockRead) Set(f func(p []byte) (r int, r1 error)) *ReadCloserMock {
	m.mock.ReadFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Read implements io.ReadCloser interface
func (m *ReadCloserMock) Read(p []byte) (r int, r1 error) {
	atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if m.ReadMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ReadMock.mockExpectations, ReadCloserMockReadParams{p},
			"ReadCloser.Read got unexpected parameters")

		if m.ReadFunc == nil {

			m.t.Fatal("No results are set for the ReadCloserMock.Read")

			return
		}
	}

	if m.ReadFunc == nil {
		m.t.Fatal("Unexpected call to ReadCloserMock.Read")
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of ReadCloserMock.ReadFunc invocations
func (m *ReadCloserMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of ReadCloserMock.Read invocations
func (m *ReadCloserMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

func (m *ReadCloserMock) minimockCheckReadCalled(f func(...interface{})) {
	if m.ReadFunc != nil && atomic.LoadUint64(&m.ReadCounter) == 0 {
		f("Expected call to ReadCloserMock.Read")
	}
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ReadCloserMock) MinimockFinish() {
	m.minimockCheck(m.t.Fatal)
}

func (m *ReadCloserMock) minimockCheck(f func(...interface{})) {
	m.minimockCheckCloseCalled(f)
	m.minimockCheckReadCalled(f)

}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ReadCloserMock) MinimockWait(timeout time.Duration) {
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
func (m *ReadCloserMock) AllMocksCalled() bool {
	return ((m.CloseFunc == nil || atomic.LoadUint64(&m.CloseCounter) > 0) &&
		(m.ReadFunc == nil || atomic.LoadUint64(&m.ReadCounter) > 0) &&
		true)
}
