package gowrap

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Writer" can be found in io
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//WriterMock implements io.Writer
type WriterMock struct {
	t minimock.Tester

	WriteFunc       func(p []byte) (r int, r1 error)
	WriteCounter    uint64
	WritePreCounter uint64
	WriteMock       mWriterMockWrite
}

//NewWriterMock returns a mock for io.Writer
func NewWriterMock(t minimock.Tester) *WriterMock {
	m := &WriterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.WriteMock = mWriterMockWrite{mock: m}

	return m
}

type mWriterMockWrite struct {
	mock             *WriterMock
	mockExpectations *WriterMockWriteParams
}

//WriterMockWriteParams represents input parameters of the Writer.Write
type WriterMockWriteParams struct {
	p []byte
}

//Expect sets up expected params for the Writer.Write
func (m *mWriterMockWrite) Expect(p []byte) *mWriterMockWrite {
	m.mockExpectations = &WriterMockWriteParams{p}
	return m
}

//Return sets up a mock for Writer.Write to return Return's arguments
func (m *mWriterMockWrite) Return(r int, r1 error) *WriterMock {
	m.mock.WriteFunc = func(p []byte) (int, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Writer.Write method
func (m *mWriterMockWrite) Set(f func(p []byte) (r int, r1 error)) *WriterMock {
	m.mock.WriteFunc = f
	return m.mock
}

//Write implements io.Writer interface
func (m *WriterMock) Write(p []byte) (r int, r1 error) {
	atomic.AddUint64(&m.WritePreCounter, 1)
	defer atomic.AddUint64(&m.WriteCounter, 1)

	if m.WriteMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.WriteMock.mockExpectations, WriterMockWriteParams{p},
			"Writer.Write got unexpected parameters")

		if m.WriteFunc == nil {

			m.t.Fatal("No results are set for the WriterMock.Write")

			return
		}
	}

	if m.WriteFunc == nil {
		m.t.Fatal("Unexpected call to WriterMock.Write")
		return
	}

	return m.WriteFunc(p)
}

//WriteMinimockCounter returns a count of WriterMock.WriteFunc invocations
func (m *WriterMock) WriteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteCounter)
}

//WriteMinimockPreCounter returns the value of WriterMock.Write invocations
func (m *WriterMock) WriteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WritePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *WriterMock) ValidateCallCounters() {

	if m.WriteFunc != nil && atomic.LoadUint64(&m.WriteCounter) == 0 {
		m.t.Fatal("Expected call to WriterMock.Write")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *WriterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *WriterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *WriterMock) MinimockFinish() {

	if m.WriteFunc != nil && atomic.LoadUint64(&m.WriteCounter) == 0 {
		m.t.Fatal("Expected call to WriterMock.Write")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *WriterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *WriterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.WriteFunc == nil || atomic.LoadUint64(&m.WriteCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.WriteFunc != nil && atomic.LoadUint64(&m.WriteCounter) == 0 {
				m.t.Error("Expected call to WriterMock.Write")
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
func (m *WriterMock) AllMocksCalled() bool {

	if m.WriteFunc != nil && atomic.LoadUint64(&m.WriteCounter) == 0 {
		return false
	}

	return true
}
