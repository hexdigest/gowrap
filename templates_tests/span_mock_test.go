package templatestests

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Span" can be found in github.com/opentracing/opentracing-go
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/opentracing/opentracing-go/log"

	testify_assert "github.com/stretchr/testify/assert"
)

//SpanMock implements github.com/opentracing/opentracing-go.Span
type SpanMock struct {
	t minimock.Tester

	BaggageItemFunc       func(p string) (r string)
	BaggageItemCounter    uint64
	BaggageItemPreCounter uint64
	BaggageItemMock       mSpanMockBaggageItem

	ContextFunc       func() (r opentracing.SpanContext)
	ContextCounter    uint64
	ContextPreCounter uint64
	ContextMock       mSpanMockContext

	FinishFunc       func()
	FinishCounter    uint64
	FinishPreCounter uint64
	FinishMock       mSpanMockFinish

	FinishWithOptionsFunc       func(p opentracing.FinishOptions)
	FinishWithOptionsCounter    uint64
	FinishWithOptionsPreCounter uint64
	FinishWithOptionsMock       mSpanMockFinishWithOptions

	LogFunc       func(p opentracing.LogData)
	LogCounter    uint64
	LogPreCounter uint64
	LogMock       mSpanMockLog

	LogEventFunc       func(p string)
	LogEventCounter    uint64
	LogEventPreCounter uint64
	LogEventMock       mSpanMockLogEvent

	LogEventWithPayloadFunc       func(p string, p1 interface{})
	LogEventWithPayloadCounter    uint64
	LogEventWithPayloadPreCounter uint64
	LogEventWithPayloadMock       mSpanMockLogEventWithPayload

	LogFieldsFunc       func(p ...log.Field)
	LogFieldsCounter    uint64
	LogFieldsPreCounter uint64
	LogFieldsMock       mSpanMockLogFields

	LogKVFunc       func(p ...interface{})
	LogKVCounter    uint64
	LogKVPreCounter uint64
	LogKVMock       mSpanMockLogKV

	SetBaggageItemFunc       func(p string, p1 string) (r opentracing.Span)
	SetBaggageItemCounter    uint64
	SetBaggageItemPreCounter uint64
	SetBaggageItemMock       mSpanMockSetBaggageItem

	SetOperationNameFunc       func(p string) (r opentracing.Span)
	SetOperationNameCounter    uint64
	SetOperationNamePreCounter uint64
	SetOperationNameMock       mSpanMockSetOperationName

	SetTagFunc       func(p string, p1 interface{}) (r opentracing.Span)
	SetTagCounter    uint64
	SetTagPreCounter uint64
	SetTagMock       mSpanMockSetTag

	TracerFunc       func() (r opentracing.Tracer)
	TracerCounter    uint64
	TracerPreCounter uint64
	TracerMock       mSpanMockTracer
}

//NewSpanMock returns a mock for github.com/opentracing/opentracing-go.Span
func NewSpanMock(t minimock.Tester) *SpanMock {
	m := &SpanMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BaggageItemMock = mSpanMockBaggageItem{mock: m}
	m.ContextMock = mSpanMockContext{mock: m}
	m.FinishMock = mSpanMockFinish{mock: m}
	m.FinishWithOptionsMock = mSpanMockFinishWithOptions{mock: m}
	m.LogMock = mSpanMockLog{mock: m}
	m.LogEventMock = mSpanMockLogEvent{mock: m}
	m.LogEventWithPayloadMock = mSpanMockLogEventWithPayload{mock: m}
	m.LogFieldsMock = mSpanMockLogFields{mock: m}
	m.LogKVMock = mSpanMockLogKV{mock: m}
	m.SetBaggageItemMock = mSpanMockSetBaggageItem{mock: m}
	m.SetOperationNameMock = mSpanMockSetOperationName{mock: m}
	m.SetTagMock = mSpanMockSetTag{mock: m}
	m.TracerMock = mSpanMockTracer{mock: m}

	return m
}

type mSpanMockBaggageItem struct {
	mock             *SpanMock
	mockExpectations *SpanMockBaggageItemParams
}

//SpanMockBaggageItemParams represents input parameters of the Span.BaggageItem
type SpanMockBaggageItemParams struct {
	p string
}

//Expect sets up expected params for the Span.BaggageItem
func (m *mSpanMockBaggageItem) Expect(p string) *mSpanMockBaggageItem {
	m.mockExpectations = &SpanMockBaggageItemParams{p}
	return m
}

//Return sets up a mock for Span.BaggageItem to return Return's arguments
func (m *mSpanMockBaggageItem) Return(r string) *SpanMock {
	m.mock.BaggageItemFunc = func(p string) string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Span.BaggageItem method
func (m *mSpanMockBaggageItem) Set(f func(p string) (r string)) *SpanMock {
	m.mock.BaggageItemFunc = f
	return m.mock
}

//BaggageItem implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) BaggageItem(p string) (r string) {
	atomic.AddUint64(&m.BaggageItemPreCounter, 1)
	defer atomic.AddUint64(&m.BaggageItemCounter, 1)

	if m.BaggageItemMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.BaggageItemMock.mockExpectations, SpanMockBaggageItemParams{p},
			"Span.BaggageItem got unexpected parameters")

		if m.BaggageItemFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.BaggageItem")

			return
		}
	}

	if m.BaggageItemFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.BaggageItem")
		return
	}

	return m.BaggageItemFunc(p)
}

//BaggageItemMinimockCounter returns a count of SpanMock.BaggageItemFunc invocations
func (m *SpanMock) BaggageItemMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BaggageItemCounter)
}

//BaggageItemMinimockPreCounter returns the value of SpanMock.BaggageItem invocations
func (m *SpanMock) BaggageItemMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BaggageItemPreCounter)
}

type mSpanMockContext struct {
	mock *SpanMock
}

//Return sets up a mock for Span.Context to return Return's arguments
func (m *mSpanMockContext) Return(r opentracing.SpanContext) *SpanMock {
	m.mock.ContextFunc = func() opentracing.SpanContext {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Span.Context method
func (m *mSpanMockContext) Set(f func() (r opentracing.SpanContext)) *SpanMock {
	m.mock.ContextFunc = f
	return m.mock
}

//Context implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) Context() (r opentracing.SpanContext) {
	atomic.AddUint64(&m.ContextPreCounter, 1)
	defer atomic.AddUint64(&m.ContextCounter, 1)

	if m.ContextFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.Context")
		return
	}

	return m.ContextFunc()
}

//ContextMinimockCounter returns a count of SpanMock.ContextFunc invocations
func (m *SpanMock) ContextMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ContextCounter)
}

//ContextMinimockPreCounter returns the value of SpanMock.Context invocations
func (m *SpanMock) ContextMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ContextPreCounter)
}

type mSpanMockFinish struct {
	mock *SpanMock
}

//Return sets up a mock for Span.Finish to return Return's arguments
func (m *mSpanMockFinish) Return() *SpanMock {
	m.mock.FinishFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Span.Finish method
func (m *mSpanMockFinish) Set(f func()) *SpanMock {
	m.mock.FinishFunc = f
	return m.mock
}

//Finish implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) Finish() {
	atomic.AddUint64(&m.FinishPreCounter, 1)
	defer atomic.AddUint64(&m.FinishCounter, 1)

	if m.FinishFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.Finish")
		return
	}

	m.FinishFunc()
}

//FinishMinimockCounter returns a count of SpanMock.FinishFunc invocations
func (m *SpanMock) FinishMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FinishCounter)
}

//FinishMinimockPreCounter returns the value of SpanMock.Finish invocations
func (m *SpanMock) FinishMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FinishPreCounter)
}

type mSpanMockFinishWithOptions struct {
	mock             *SpanMock
	mockExpectations *SpanMockFinishWithOptionsParams
}

//SpanMockFinishWithOptionsParams represents input parameters of the Span.FinishWithOptions
type SpanMockFinishWithOptionsParams struct {
	p opentracing.FinishOptions
}

//Expect sets up expected params for the Span.FinishWithOptions
func (m *mSpanMockFinishWithOptions) Expect(p opentracing.FinishOptions) *mSpanMockFinishWithOptions {
	m.mockExpectations = &SpanMockFinishWithOptionsParams{p}
	return m
}

//Return sets up a mock for Span.FinishWithOptions to return Return's arguments
func (m *mSpanMockFinishWithOptions) Return() *SpanMock {
	m.mock.FinishWithOptionsFunc = func(p opentracing.FinishOptions) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Span.FinishWithOptions method
func (m *mSpanMockFinishWithOptions) Set(f func(p opentracing.FinishOptions)) *SpanMock {
	m.mock.FinishWithOptionsFunc = f
	return m.mock
}

//FinishWithOptions implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) FinishWithOptions(p opentracing.FinishOptions) {
	atomic.AddUint64(&m.FinishWithOptionsPreCounter, 1)
	defer atomic.AddUint64(&m.FinishWithOptionsCounter, 1)

	if m.FinishWithOptionsMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.FinishWithOptionsMock.mockExpectations, SpanMockFinishWithOptionsParams{p},
			"Span.FinishWithOptions got unexpected parameters")

		if m.FinishWithOptionsFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.FinishWithOptions")

			return
		}
	}

	if m.FinishWithOptionsFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.FinishWithOptions")
		return
	}

	m.FinishWithOptionsFunc(p)
}

//FinishWithOptionsMinimockCounter returns a count of SpanMock.FinishWithOptionsFunc invocations
func (m *SpanMock) FinishWithOptionsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FinishWithOptionsCounter)
}

//FinishWithOptionsMinimockPreCounter returns the value of SpanMock.FinishWithOptions invocations
func (m *SpanMock) FinishWithOptionsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FinishWithOptionsPreCounter)
}

type mSpanMockLog struct {
	mock             *SpanMock
	mockExpectations *SpanMockLogParams
}

//SpanMockLogParams represents input parameters of the Span.Log
type SpanMockLogParams struct {
	p opentracing.LogData
}

//Expect sets up expected params for the Span.Log
func (m *mSpanMockLog) Expect(p opentracing.LogData) *mSpanMockLog {
	m.mockExpectations = &SpanMockLogParams{p}
	return m
}

//Return sets up a mock for Span.Log to return Return's arguments
func (m *mSpanMockLog) Return() *SpanMock {
	m.mock.LogFunc = func(p opentracing.LogData) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Span.Log method
func (m *mSpanMockLog) Set(f func(p opentracing.LogData)) *SpanMock {
	m.mock.LogFunc = f
	return m.mock
}

//Log implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) Log(p opentracing.LogData) {
	atomic.AddUint64(&m.LogPreCounter, 1)
	defer atomic.AddUint64(&m.LogCounter, 1)

	if m.LogMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LogMock.mockExpectations, SpanMockLogParams{p},
			"Span.Log got unexpected parameters")

		if m.LogFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.Log")

			return
		}
	}

	if m.LogFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.Log")
		return
	}

	m.LogFunc(p)
}

//LogMinimockCounter returns a count of SpanMock.LogFunc invocations
func (m *SpanMock) LogMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LogCounter)
}

//LogMinimockPreCounter returns the value of SpanMock.Log invocations
func (m *SpanMock) LogMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LogPreCounter)
}

type mSpanMockLogEvent struct {
	mock             *SpanMock
	mockExpectations *SpanMockLogEventParams
}

//SpanMockLogEventParams represents input parameters of the Span.LogEvent
type SpanMockLogEventParams struct {
	p string
}

//Expect sets up expected params for the Span.LogEvent
func (m *mSpanMockLogEvent) Expect(p string) *mSpanMockLogEvent {
	m.mockExpectations = &SpanMockLogEventParams{p}
	return m
}

//Return sets up a mock for Span.LogEvent to return Return's arguments
func (m *mSpanMockLogEvent) Return() *SpanMock {
	m.mock.LogEventFunc = func(p string) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Span.LogEvent method
func (m *mSpanMockLogEvent) Set(f func(p string)) *SpanMock {
	m.mock.LogEventFunc = f
	return m.mock
}

//LogEvent implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) LogEvent(p string) {
	atomic.AddUint64(&m.LogEventPreCounter, 1)
	defer atomic.AddUint64(&m.LogEventCounter, 1)

	if m.LogEventMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LogEventMock.mockExpectations, SpanMockLogEventParams{p},
			"Span.LogEvent got unexpected parameters")

		if m.LogEventFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.LogEvent")

			return
		}
	}

	if m.LogEventFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.LogEvent")
		return
	}

	m.LogEventFunc(p)
}

//LogEventMinimockCounter returns a count of SpanMock.LogEventFunc invocations
func (m *SpanMock) LogEventMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LogEventCounter)
}

//LogEventMinimockPreCounter returns the value of SpanMock.LogEvent invocations
func (m *SpanMock) LogEventMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LogEventPreCounter)
}

type mSpanMockLogEventWithPayload struct {
	mock             *SpanMock
	mockExpectations *SpanMockLogEventWithPayloadParams
}

//SpanMockLogEventWithPayloadParams represents input parameters of the Span.LogEventWithPayload
type SpanMockLogEventWithPayloadParams struct {
	p  string
	p1 interface{}
}

//Expect sets up expected params for the Span.LogEventWithPayload
func (m *mSpanMockLogEventWithPayload) Expect(p string, p1 interface{}) *mSpanMockLogEventWithPayload {
	m.mockExpectations = &SpanMockLogEventWithPayloadParams{p, p1}
	return m
}

//Return sets up a mock for Span.LogEventWithPayload to return Return's arguments
func (m *mSpanMockLogEventWithPayload) Return() *SpanMock {
	m.mock.LogEventWithPayloadFunc = func(p string, p1 interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Span.LogEventWithPayload method
func (m *mSpanMockLogEventWithPayload) Set(f func(p string, p1 interface{})) *SpanMock {
	m.mock.LogEventWithPayloadFunc = f
	return m.mock
}

//LogEventWithPayload implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) LogEventWithPayload(p string, p1 interface{}) {
	atomic.AddUint64(&m.LogEventWithPayloadPreCounter, 1)
	defer atomic.AddUint64(&m.LogEventWithPayloadCounter, 1)

	if m.LogEventWithPayloadMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LogEventWithPayloadMock.mockExpectations, SpanMockLogEventWithPayloadParams{p, p1},
			"Span.LogEventWithPayload got unexpected parameters")

		if m.LogEventWithPayloadFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.LogEventWithPayload")

			return
		}
	}

	if m.LogEventWithPayloadFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.LogEventWithPayload")
		return
	}

	m.LogEventWithPayloadFunc(p, p1)
}

//LogEventWithPayloadMinimockCounter returns a count of SpanMock.LogEventWithPayloadFunc invocations
func (m *SpanMock) LogEventWithPayloadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LogEventWithPayloadCounter)
}

//LogEventWithPayloadMinimockPreCounter returns the value of SpanMock.LogEventWithPayload invocations
func (m *SpanMock) LogEventWithPayloadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LogEventWithPayloadPreCounter)
}

type mSpanMockLogFields struct {
	mock             *SpanMock
	mockExpectations *SpanMockLogFieldsParams
}

//SpanMockLogFieldsParams represents input parameters of the Span.LogFields
type SpanMockLogFieldsParams struct {
	p []log.Field
}

//Expect sets up expected params for the Span.LogFields
func (m *mSpanMockLogFields) Expect(p ...log.Field) *mSpanMockLogFields {
	m.mockExpectations = &SpanMockLogFieldsParams{p}
	return m
}

//Return sets up a mock for Span.LogFields to return Return's arguments
func (m *mSpanMockLogFields) Return() *SpanMock {
	m.mock.LogFieldsFunc = func(p ...log.Field) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Span.LogFields method
func (m *mSpanMockLogFields) Set(f func(p ...log.Field)) *SpanMock {
	m.mock.LogFieldsFunc = f
	return m.mock
}

//LogFields implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) LogFields(p ...log.Field) {
	atomic.AddUint64(&m.LogFieldsPreCounter, 1)
	defer atomic.AddUint64(&m.LogFieldsCounter, 1)

	if m.LogFieldsMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LogFieldsMock.mockExpectations, SpanMockLogFieldsParams{p},
			"Span.LogFields got unexpected parameters")

		if m.LogFieldsFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.LogFields")

			return
		}
	}

	if m.LogFieldsFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.LogFields")
		return
	}

	m.LogFieldsFunc(p...)
}

//LogFieldsMinimockCounter returns a count of SpanMock.LogFieldsFunc invocations
func (m *SpanMock) LogFieldsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LogFieldsCounter)
}

//LogFieldsMinimockPreCounter returns the value of SpanMock.LogFields invocations
func (m *SpanMock) LogFieldsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LogFieldsPreCounter)
}

type mSpanMockLogKV struct {
	mock             *SpanMock
	mockExpectations *SpanMockLogKVParams
}

//SpanMockLogKVParams represents input parameters of the Span.LogKV
type SpanMockLogKVParams struct {
	p []interface{}
}

//Expect sets up expected params for the Span.LogKV
func (m *mSpanMockLogKV) Expect(p ...interface{}) *mSpanMockLogKV {
	m.mockExpectations = &SpanMockLogKVParams{p}
	return m
}

//Return sets up a mock for Span.LogKV to return Return's arguments
func (m *mSpanMockLogKV) Return() *SpanMock {
	m.mock.LogKVFunc = func(p ...interface{}) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Span.LogKV method
func (m *mSpanMockLogKV) Set(f func(p ...interface{})) *SpanMock {
	m.mock.LogKVFunc = f
	return m.mock
}

//LogKV implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) LogKV(p ...interface{}) {
	atomic.AddUint64(&m.LogKVPreCounter, 1)
	defer atomic.AddUint64(&m.LogKVCounter, 1)

	if m.LogKVMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LogKVMock.mockExpectations, SpanMockLogKVParams{p},
			"Span.LogKV got unexpected parameters")

		if m.LogKVFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.LogKV")

			return
		}
	}

	if m.LogKVFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.LogKV")
		return
	}

	m.LogKVFunc(p...)
}

//LogKVMinimockCounter returns a count of SpanMock.LogKVFunc invocations
func (m *SpanMock) LogKVMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LogKVCounter)
}

//LogKVMinimockPreCounter returns the value of SpanMock.LogKV invocations
func (m *SpanMock) LogKVMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LogKVPreCounter)
}

type mSpanMockSetBaggageItem struct {
	mock             *SpanMock
	mockExpectations *SpanMockSetBaggageItemParams
}

//SpanMockSetBaggageItemParams represents input parameters of the Span.SetBaggageItem
type SpanMockSetBaggageItemParams struct {
	p  string
	p1 string
}

//Expect sets up expected params for the Span.SetBaggageItem
func (m *mSpanMockSetBaggageItem) Expect(p string, p1 string) *mSpanMockSetBaggageItem {
	m.mockExpectations = &SpanMockSetBaggageItemParams{p, p1}
	return m
}

//Return sets up a mock for Span.SetBaggageItem to return Return's arguments
func (m *mSpanMockSetBaggageItem) Return(r opentracing.Span) *SpanMock {
	m.mock.SetBaggageItemFunc = func(p string, p1 string) opentracing.Span {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Span.SetBaggageItem method
func (m *mSpanMockSetBaggageItem) Set(f func(p string, p1 string) (r opentracing.Span)) *SpanMock {
	m.mock.SetBaggageItemFunc = f
	return m.mock
}

//SetBaggageItem implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) SetBaggageItem(p string, p1 string) (r opentracing.Span) {
	atomic.AddUint64(&m.SetBaggageItemPreCounter, 1)
	defer atomic.AddUint64(&m.SetBaggageItemCounter, 1)

	if m.SetBaggageItemMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetBaggageItemMock.mockExpectations, SpanMockSetBaggageItemParams{p, p1},
			"Span.SetBaggageItem got unexpected parameters")

		if m.SetBaggageItemFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.SetBaggageItem")

			return
		}
	}

	if m.SetBaggageItemFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.SetBaggageItem")
		return
	}

	return m.SetBaggageItemFunc(p, p1)
}

//SetBaggageItemMinimockCounter returns a count of SpanMock.SetBaggageItemFunc invocations
func (m *SpanMock) SetBaggageItemMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetBaggageItemCounter)
}

//SetBaggageItemMinimockPreCounter returns the value of SpanMock.SetBaggageItem invocations
func (m *SpanMock) SetBaggageItemMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetBaggageItemPreCounter)
}

type mSpanMockSetOperationName struct {
	mock             *SpanMock
	mockExpectations *SpanMockSetOperationNameParams
}

//SpanMockSetOperationNameParams represents input parameters of the Span.SetOperationName
type SpanMockSetOperationNameParams struct {
	p string
}

//Expect sets up expected params for the Span.SetOperationName
func (m *mSpanMockSetOperationName) Expect(p string) *mSpanMockSetOperationName {
	m.mockExpectations = &SpanMockSetOperationNameParams{p}
	return m
}

//Return sets up a mock for Span.SetOperationName to return Return's arguments
func (m *mSpanMockSetOperationName) Return(r opentracing.Span) *SpanMock {
	m.mock.SetOperationNameFunc = func(p string) opentracing.Span {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Span.SetOperationName method
func (m *mSpanMockSetOperationName) Set(f func(p string) (r opentracing.Span)) *SpanMock {
	m.mock.SetOperationNameFunc = f
	return m.mock
}

//SetOperationName implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) SetOperationName(p string) (r opentracing.Span) {
	atomic.AddUint64(&m.SetOperationNamePreCounter, 1)
	defer atomic.AddUint64(&m.SetOperationNameCounter, 1)

	if m.SetOperationNameMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetOperationNameMock.mockExpectations, SpanMockSetOperationNameParams{p},
			"Span.SetOperationName got unexpected parameters")

		if m.SetOperationNameFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.SetOperationName")

			return
		}
	}

	if m.SetOperationNameFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.SetOperationName")
		return
	}

	return m.SetOperationNameFunc(p)
}

//SetOperationNameMinimockCounter returns a count of SpanMock.SetOperationNameFunc invocations
func (m *SpanMock) SetOperationNameMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetOperationNameCounter)
}

//SetOperationNameMinimockPreCounter returns the value of SpanMock.SetOperationName invocations
func (m *SpanMock) SetOperationNameMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetOperationNamePreCounter)
}

type mSpanMockSetTag struct {
	mock             *SpanMock
	mockExpectations *SpanMockSetTagParams
}

//SpanMockSetTagParams represents input parameters of the Span.SetTag
type SpanMockSetTagParams struct {
	p  string
	p1 interface{}
}

//Expect sets up expected params for the Span.SetTag
func (m *mSpanMockSetTag) Expect(p string, p1 interface{}) *mSpanMockSetTag {
	m.mockExpectations = &SpanMockSetTagParams{p, p1}
	return m
}

//Return sets up a mock for Span.SetTag to return Return's arguments
func (m *mSpanMockSetTag) Return(r opentracing.Span) *SpanMock {
	m.mock.SetTagFunc = func(p string, p1 interface{}) opentracing.Span {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Span.SetTag method
func (m *mSpanMockSetTag) Set(f func(p string, p1 interface{}) (r opentracing.Span)) *SpanMock {
	m.mock.SetTagFunc = f
	return m.mock
}

//SetTag implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) SetTag(p string, p1 interface{}) (r opentracing.Span) {
	atomic.AddUint64(&m.SetTagPreCounter, 1)
	defer atomic.AddUint64(&m.SetTagCounter, 1)

	if m.SetTagMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetTagMock.mockExpectations, SpanMockSetTagParams{p, p1},
			"Span.SetTag got unexpected parameters")

		if m.SetTagFunc == nil {

			m.t.Fatal("No results are set for the SpanMock.SetTag")

			return
		}
	}

	if m.SetTagFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.SetTag")
		return
	}

	return m.SetTagFunc(p, p1)
}

//SetTagMinimockCounter returns a count of SpanMock.SetTagFunc invocations
func (m *SpanMock) SetTagMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetTagCounter)
}

//SetTagMinimockPreCounter returns the value of SpanMock.SetTag invocations
func (m *SpanMock) SetTagMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetTagPreCounter)
}

type mSpanMockTracer struct {
	mock *SpanMock
}

//Return sets up a mock for Span.Tracer to return Return's arguments
func (m *mSpanMockTracer) Return(r opentracing.Tracer) *SpanMock {
	m.mock.TracerFunc = func() opentracing.Tracer {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Span.Tracer method
func (m *mSpanMockTracer) Set(f func() (r opentracing.Tracer)) *SpanMock {
	m.mock.TracerFunc = f
	return m.mock
}

//Tracer implements github.com/opentracing/opentracing-go.Span interface
func (m *SpanMock) Tracer() (r opentracing.Tracer) {
	atomic.AddUint64(&m.TracerPreCounter, 1)
	defer atomic.AddUint64(&m.TracerCounter, 1)

	if m.TracerFunc == nil {
		m.t.Fatal("Unexpected call to SpanMock.Tracer")
		return
	}

	return m.TracerFunc()
}

//TracerMinimockCounter returns a count of SpanMock.TracerFunc invocations
func (m *SpanMock) TracerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.TracerCounter)
}

//TracerMinimockPreCounter returns the value of SpanMock.Tracer invocations
func (m *SpanMock) TracerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.TracerPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SpanMock) ValidateCallCounters() {

	if m.BaggageItemFunc != nil && atomic.LoadUint64(&m.BaggageItemCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.BaggageItem")
	}

	if m.ContextFunc != nil && atomic.LoadUint64(&m.ContextCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.Context")
	}

	if m.FinishFunc != nil && atomic.LoadUint64(&m.FinishCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.Finish")
	}

	if m.FinishWithOptionsFunc != nil && atomic.LoadUint64(&m.FinishWithOptionsCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.FinishWithOptions")
	}

	if m.LogFunc != nil && atomic.LoadUint64(&m.LogCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.Log")
	}

	if m.LogEventFunc != nil && atomic.LoadUint64(&m.LogEventCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.LogEvent")
	}

	if m.LogEventWithPayloadFunc != nil && atomic.LoadUint64(&m.LogEventWithPayloadCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.LogEventWithPayload")
	}

	if m.LogFieldsFunc != nil && atomic.LoadUint64(&m.LogFieldsCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.LogFields")
	}

	if m.LogKVFunc != nil && atomic.LoadUint64(&m.LogKVCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.LogKV")
	}

	if m.SetBaggageItemFunc != nil && atomic.LoadUint64(&m.SetBaggageItemCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.SetBaggageItem")
	}

	if m.SetOperationNameFunc != nil && atomic.LoadUint64(&m.SetOperationNameCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.SetOperationName")
	}

	if m.SetTagFunc != nil && atomic.LoadUint64(&m.SetTagCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.SetTag")
	}

	if m.TracerFunc != nil && atomic.LoadUint64(&m.TracerCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.Tracer")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SpanMock) CheckMocksCalled() {
	m.Finish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SpanMock) MinimockFinish() {

	if m.BaggageItemFunc != nil && atomic.LoadUint64(&m.BaggageItemCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.BaggageItem")
	}

	if m.ContextFunc != nil && atomic.LoadUint64(&m.ContextCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.Context")
	}

	if m.FinishFunc != nil && atomic.LoadUint64(&m.FinishCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.Finish")
	}

	if m.FinishWithOptionsFunc != nil && atomic.LoadUint64(&m.FinishWithOptionsCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.FinishWithOptions")
	}

	if m.LogFunc != nil && atomic.LoadUint64(&m.LogCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.Log")
	}

	if m.LogEventFunc != nil && atomic.LoadUint64(&m.LogEventCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.LogEvent")
	}

	if m.LogEventWithPayloadFunc != nil && atomic.LoadUint64(&m.LogEventWithPayloadCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.LogEventWithPayload")
	}

	if m.LogFieldsFunc != nil && atomic.LoadUint64(&m.LogFieldsCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.LogFields")
	}

	if m.LogKVFunc != nil && atomic.LoadUint64(&m.LogKVCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.LogKV")
	}

	if m.SetBaggageItemFunc != nil && atomic.LoadUint64(&m.SetBaggageItemCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.SetBaggageItem")
	}

	if m.SetOperationNameFunc != nil && atomic.LoadUint64(&m.SetOperationNameCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.SetOperationName")
	}

	if m.SetTagFunc != nil && atomic.LoadUint64(&m.SetTagCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.SetTag")
	}

	if m.TracerFunc != nil && atomic.LoadUint64(&m.TracerCounter) == 0 {
		m.t.Fatal("Expected call to SpanMock.Tracer")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SpanMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SpanMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.BaggageItemFunc == nil || atomic.LoadUint64(&m.BaggageItemCounter) > 0)
		ok = ok && (m.ContextFunc == nil || atomic.LoadUint64(&m.ContextCounter) > 0)
		ok = ok && (m.FinishFunc == nil || atomic.LoadUint64(&m.FinishCounter) > 0)
		ok = ok && (m.FinishWithOptionsFunc == nil || atomic.LoadUint64(&m.FinishWithOptionsCounter) > 0)
		ok = ok && (m.LogFunc == nil || atomic.LoadUint64(&m.LogCounter) > 0)
		ok = ok && (m.LogEventFunc == nil || atomic.LoadUint64(&m.LogEventCounter) > 0)
		ok = ok && (m.LogEventWithPayloadFunc == nil || atomic.LoadUint64(&m.LogEventWithPayloadCounter) > 0)
		ok = ok && (m.LogFieldsFunc == nil || atomic.LoadUint64(&m.LogFieldsCounter) > 0)
		ok = ok && (m.LogKVFunc == nil || atomic.LoadUint64(&m.LogKVCounter) > 0)
		ok = ok && (m.SetBaggageItemFunc == nil || atomic.LoadUint64(&m.SetBaggageItemCounter) > 0)
		ok = ok && (m.SetOperationNameFunc == nil || atomic.LoadUint64(&m.SetOperationNameCounter) > 0)
		ok = ok && (m.SetTagFunc == nil || atomic.LoadUint64(&m.SetTagCounter) > 0)
		ok = ok && (m.TracerFunc == nil || atomic.LoadUint64(&m.TracerCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.BaggageItemFunc != nil && atomic.LoadUint64(&m.BaggageItemCounter) == 0 {
				m.t.Error("Expected call to SpanMock.BaggageItem")
			}

			if m.ContextFunc != nil && atomic.LoadUint64(&m.ContextCounter) == 0 {
				m.t.Error("Expected call to SpanMock.Context")
			}

			if m.FinishFunc != nil && atomic.LoadUint64(&m.FinishCounter) == 0 {
				m.t.Error("Expected call to SpanMock.Finish")
			}

			if m.FinishWithOptionsFunc != nil && atomic.LoadUint64(&m.FinishWithOptionsCounter) == 0 {
				m.t.Error("Expected call to SpanMock.FinishWithOptions")
			}

			if m.LogFunc != nil && atomic.LoadUint64(&m.LogCounter) == 0 {
				m.t.Error("Expected call to SpanMock.Log")
			}

			if m.LogEventFunc != nil && atomic.LoadUint64(&m.LogEventCounter) == 0 {
				m.t.Error("Expected call to SpanMock.LogEvent")
			}

			if m.LogEventWithPayloadFunc != nil && atomic.LoadUint64(&m.LogEventWithPayloadCounter) == 0 {
				m.t.Error("Expected call to SpanMock.LogEventWithPayload")
			}

			if m.LogFieldsFunc != nil && atomic.LoadUint64(&m.LogFieldsCounter) == 0 {
				m.t.Error("Expected call to SpanMock.LogFields")
			}

			if m.LogKVFunc != nil && atomic.LoadUint64(&m.LogKVCounter) == 0 {
				m.t.Error("Expected call to SpanMock.LogKV")
			}

			if m.SetBaggageItemFunc != nil && atomic.LoadUint64(&m.SetBaggageItemCounter) == 0 {
				m.t.Error("Expected call to SpanMock.SetBaggageItem")
			}

			if m.SetOperationNameFunc != nil && atomic.LoadUint64(&m.SetOperationNameCounter) == 0 {
				m.t.Error("Expected call to SpanMock.SetOperationName")
			}

			if m.SetTagFunc != nil && atomic.LoadUint64(&m.SetTagCounter) == 0 {
				m.t.Error("Expected call to SpanMock.SetTag")
			}

			if m.TracerFunc != nil && atomic.LoadUint64(&m.TracerCounter) == 0 {
				m.t.Error("Expected call to SpanMock.Tracer")
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
func (m *SpanMock) AllMocksCalled() bool {

	if m.BaggageItemFunc != nil && atomic.LoadUint64(&m.BaggageItemCounter) == 0 {
		return false
	}

	if m.ContextFunc != nil && atomic.LoadUint64(&m.ContextCounter) == 0 {
		return false
	}

	if m.FinishFunc != nil && atomic.LoadUint64(&m.FinishCounter) == 0 {
		return false
	}

	if m.FinishWithOptionsFunc != nil && atomic.LoadUint64(&m.FinishWithOptionsCounter) == 0 {
		return false
	}

	if m.LogFunc != nil && atomic.LoadUint64(&m.LogCounter) == 0 {
		return false
	}

	if m.LogEventFunc != nil && atomic.LoadUint64(&m.LogEventCounter) == 0 {
		return false
	}

	if m.LogEventWithPayloadFunc != nil && atomic.LoadUint64(&m.LogEventWithPayloadCounter) == 0 {
		return false
	}

	if m.LogFieldsFunc != nil && atomic.LoadUint64(&m.LogFieldsCounter) == 0 {
		return false
	}

	if m.LogKVFunc != nil && atomic.LoadUint64(&m.LogKVCounter) == 0 {
		return false
	}

	if m.SetBaggageItemFunc != nil && atomic.LoadUint64(&m.SetBaggageItemCounter) == 0 {
		return false
	}

	if m.SetOperationNameFunc != nil && atomic.LoadUint64(&m.SetOperationNameCounter) == 0 {
		return false
	}

	if m.SetTagFunc != nil && atomic.LoadUint64(&m.SetTagCounter) == 0 {
		return false
	}

	if m.TracerFunc != nil && atomic.LoadUint64(&m.TracerCounter) == 0 {
		return false
	}

	return true
}
