package pkg

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FileInfo" can be found in os
*/
import (
	os "os"
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
)

//FileInfoMock implements os.FileInfo
type FileInfoMock struct {
	t minimock.Tester

	IsDirFunc       func() (r bool)
	IsDirCounter    uint64
	IsDirPreCounter uint64
	IsDirMock       mFileInfoMockIsDir

	ModTimeFunc       func() (r time.Time)
	ModTimeCounter    uint64
	ModTimePreCounter uint64
	ModTimeMock       mFileInfoMockModTime

	ModeFunc       func() (r os.FileMode)
	ModeCounter    uint64
	ModePreCounter uint64
	ModeMock       mFileInfoMockMode

	NameFunc       func() (r string)
	NameCounter    uint64
	NamePreCounter uint64
	NameMock       mFileInfoMockName

	SizeFunc       func() (r int64)
	SizeCounter    uint64
	SizePreCounter uint64
	SizeMock       mFileInfoMockSize

	SysFunc       func() (r interface{})
	SysCounter    uint64
	SysPreCounter uint64
	SysMock       mFileInfoMockSys
}

//NewFileInfoMock returns a mock for os.FileInfo
func NewFileInfoMock(t minimock.Tester) *FileInfoMock {
	m := &FileInfoMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IsDirMock = mFileInfoMockIsDir{mock: m}
	m.ModTimeMock = mFileInfoMockModTime{mock: m}
	m.ModeMock = mFileInfoMockMode{mock: m}
	m.NameMock = mFileInfoMockName{mock: m}
	m.SizeMock = mFileInfoMockSize{mock: m}
	m.SysMock = mFileInfoMockSys{mock: m}

	return m
}

type mFileInfoMockIsDir struct {
	mock *FileInfoMock
}

//Return sets up a mock for FileInfo.IsDir to return Return's arguments
func (m *mFileInfoMockIsDir) Return(r bool) *FileInfoMock {
	m.mock.IsDirFunc = func() bool {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of FileInfo.IsDir method
func (m *mFileInfoMockIsDir) Set(f func() (r bool)) *FileInfoMock {
	m.mock.IsDirFunc = f

	return m.mock
}

//IsDir implements os.FileInfo interface
func (m *FileInfoMock) IsDir() (r bool) {
	atomic.AddUint64(&m.IsDirPreCounter, 1)
	defer atomic.AddUint64(&m.IsDirCounter, 1)

	if m.IsDirFunc == nil {
		m.t.Fatal("Unexpected call to FileInfoMock.IsDir")
		return
	}

	return m.IsDirFunc()
}

//IsDirMinimockCounter returns a count of FileInfoMock.IsDirFunc invocations
func (m *FileInfoMock) IsDirMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsDirCounter)
}

//IsDirMinimockPreCounter returns the value of FileInfoMock.IsDir invocations
func (m *FileInfoMock) IsDirMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsDirPreCounter)
}

func (m *FileInfoMock) minimockCheckIsDirCalled(f func(...interface{})) {
	if m.IsDirFunc != nil && atomic.LoadUint64(&m.IsDirCounter) == 0 {
		f("Expected call to FileInfoMock.IsDir")
	}
}

type mFileInfoMockModTime struct {
	mock *FileInfoMock
}

//Return sets up a mock for FileInfo.ModTime to return Return's arguments
func (m *mFileInfoMockModTime) Return(r time.Time) *FileInfoMock {
	m.mock.ModTimeFunc = func() time.Time {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of FileInfo.ModTime method
func (m *mFileInfoMockModTime) Set(f func() (r time.Time)) *FileInfoMock {
	m.mock.ModTimeFunc = f

	return m.mock
}

//ModTime implements os.FileInfo interface
func (m *FileInfoMock) ModTime() (r time.Time) {
	atomic.AddUint64(&m.ModTimePreCounter, 1)
	defer atomic.AddUint64(&m.ModTimeCounter, 1)

	if m.ModTimeFunc == nil {
		m.t.Fatal("Unexpected call to FileInfoMock.ModTime")
		return
	}

	return m.ModTimeFunc()
}

//ModTimeMinimockCounter returns a count of FileInfoMock.ModTimeFunc invocations
func (m *FileInfoMock) ModTimeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ModTimeCounter)
}

//ModTimeMinimockPreCounter returns the value of FileInfoMock.ModTime invocations
func (m *FileInfoMock) ModTimeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ModTimePreCounter)
}

func (m *FileInfoMock) minimockCheckModTimeCalled(f func(...interface{})) {
	if m.ModTimeFunc != nil && atomic.LoadUint64(&m.ModTimeCounter) == 0 {
		f("Expected call to FileInfoMock.ModTime")
	}
}

type mFileInfoMockMode struct {
	mock *FileInfoMock
}

//Return sets up a mock for FileInfo.Mode to return Return's arguments
func (m *mFileInfoMockMode) Return(r os.FileMode) *FileInfoMock {
	m.mock.ModeFunc = func() os.FileMode {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of FileInfo.Mode method
func (m *mFileInfoMockMode) Set(f func() (r os.FileMode)) *FileInfoMock {
	m.mock.ModeFunc = f

	return m.mock
}

//Mode implements os.FileInfo interface
func (m *FileInfoMock) Mode() (r os.FileMode) {
	atomic.AddUint64(&m.ModePreCounter, 1)
	defer atomic.AddUint64(&m.ModeCounter, 1)

	if m.ModeFunc == nil {
		m.t.Fatal("Unexpected call to FileInfoMock.Mode")
		return
	}

	return m.ModeFunc()
}

//ModeMinimockCounter returns a count of FileInfoMock.ModeFunc invocations
func (m *FileInfoMock) ModeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ModeCounter)
}

//ModeMinimockPreCounter returns the value of FileInfoMock.Mode invocations
func (m *FileInfoMock) ModeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ModePreCounter)
}

func (m *FileInfoMock) minimockCheckModeCalled(f func(...interface{})) {
	if m.ModeFunc != nil && atomic.LoadUint64(&m.ModeCounter) == 0 {
		f("Expected call to FileInfoMock.Mode")
	}
}

type mFileInfoMockName struct {
	mock *FileInfoMock
}

//Return sets up a mock for FileInfo.Name to return Return's arguments
func (m *mFileInfoMockName) Return(r string) *FileInfoMock {
	m.mock.NameFunc = func() string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of FileInfo.Name method
func (m *mFileInfoMockName) Set(f func() (r string)) *FileInfoMock {
	m.mock.NameFunc = f

	return m.mock
}

//Name implements os.FileInfo interface
func (m *FileInfoMock) Name() (r string) {
	atomic.AddUint64(&m.NamePreCounter, 1)
	defer atomic.AddUint64(&m.NameCounter, 1)

	if m.NameFunc == nil {
		m.t.Fatal("Unexpected call to FileInfoMock.Name")
		return
	}

	return m.NameFunc()
}

//NameMinimockCounter returns a count of FileInfoMock.NameFunc invocations
func (m *FileInfoMock) NameMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NameCounter)
}

//NameMinimockPreCounter returns the value of FileInfoMock.Name invocations
func (m *FileInfoMock) NameMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NamePreCounter)
}

func (m *FileInfoMock) minimockCheckNameCalled(f func(...interface{})) {
	if m.NameFunc != nil && atomic.LoadUint64(&m.NameCounter) == 0 {
		f("Expected call to FileInfoMock.Name")
	}
}

type mFileInfoMockSize struct {
	mock *FileInfoMock
}

//Return sets up a mock for FileInfo.Size to return Return's arguments
func (m *mFileInfoMockSize) Return(r int64) *FileInfoMock {
	m.mock.SizeFunc = func() int64 {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of FileInfo.Size method
func (m *mFileInfoMockSize) Set(f func() (r int64)) *FileInfoMock {
	m.mock.SizeFunc = f

	return m.mock
}

//Size implements os.FileInfo interface
func (m *FileInfoMock) Size() (r int64) {
	atomic.AddUint64(&m.SizePreCounter, 1)
	defer atomic.AddUint64(&m.SizeCounter, 1)

	if m.SizeFunc == nil {
		m.t.Fatal("Unexpected call to FileInfoMock.Size")
		return
	}

	return m.SizeFunc()
}

//SizeMinimockCounter returns a count of FileInfoMock.SizeFunc invocations
func (m *FileInfoMock) SizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SizeCounter)
}

//SizeMinimockPreCounter returns the value of FileInfoMock.Size invocations
func (m *FileInfoMock) SizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SizePreCounter)
}

func (m *FileInfoMock) minimockCheckSizeCalled(f func(...interface{})) {
	if m.SizeFunc != nil && atomic.LoadUint64(&m.SizeCounter) == 0 {
		f("Expected call to FileInfoMock.Size")
	}
}

type mFileInfoMockSys struct {
	mock *FileInfoMock
}

//Return sets up a mock for FileInfo.Sys to return Return's arguments
func (m *mFileInfoMockSys) Return(r interface{}) *FileInfoMock {
	m.mock.SysFunc = func() interface{} {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of FileInfo.Sys method
func (m *mFileInfoMockSys) Set(f func() (r interface{})) *FileInfoMock {
	m.mock.SysFunc = f

	return m.mock
}

//Sys implements os.FileInfo interface
func (m *FileInfoMock) Sys() (r interface{}) {
	atomic.AddUint64(&m.SysPreCounter, 1)
	defer atomic.AddUint64(&m.SysCounter, 1)

	if m.SysFunc == nil {
		m.t.Fatal("Unexpected call to FileInfoMock.Sys")
		return
	}

	return m.SysFunc()
}

//SysMinimockCounter returns a count of FileInfoMock.SysFunc invocations
func (m *FileInfoMock) SysMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SysCounter)
}

//SysMinimockPreCounter returns the value of FileInfoMock.Sys invocations
func (m *FileInfoMock) SysMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SysPreCounter)
}

func (m *FileInfoMock) minimockCheckSysCalled(f func(...interface{})) {
	if m.SysFunc != nil && atomic.LoadUint64(&m.SysCounter) == 0 {
		f("Expected call to FileInfoMock.Sys")
	}
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FileInfoMock) MinimockFinish() {
	m.minimockCheck(m.t.Fatal)
}

func (m *FileInfoMock) minimockCheck(f func(...interface{})) {
	m.minimockCheckIsDirCalled(f)
	m.minimockCheckModTimeCalled(f)
	m.minimockCheckModeCalled(f)
	m.minimockCheckNameCalled(f)
	m.minimockCheckSizeCalled(f)
	m.minimockCheckSysCalled(f)

}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FileInfoMock) MinimockWait(timeout time.Duration) {
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
func (m *FileInfoMock) AllMocksCalled() bool {
	return ((m.IsDirFunc == nil || atomic.LoadUint64(&m.IsDirCounter) > 0) &&
		(m.ModTimeFunc == nil || atomic.LoadUint64(&m.ModTimeCounter) > 0) &&
		(m.ModeFunc == nil || atomic.LoadUint64(&m.ModeCounter) > 0) &&
		(m.NameFunc == nil || atomic.LoadUint64(&m.NameCounter) > 0) &&
		(m.SizeFunc == nil || atomic.LoadUint64(&m.SizeCounter) > 0) &&
		(m.SysFunc == nil || atomic.LoadUint64(&m.SysCounter) > 0) &&
		true)
}
