package templatestests

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using ../templates/logrus template

//go:generate gowrap gen -d . -i TestInterface -t ../templates/logrus -o interface_with_logrus.go

import (
	"context"

	"github.com/Sirupsen/logrus"
)

// TestInterfaceWithLogrus implements TestInterface that is instrumented with logrus logger
type TestInterfaceWithLogrus struct {
	_log  *logrus.Entry
	_base TestInterface
}

// NewTestInterfaceWithLogrus instruments an implementation of the TestInterface with simple logging
func NewTestInterfaceWithLogrus(base TestInterface, log *logrus.Entry) TestInterfaceWithLogrus {
	return TestInterfaceWithLogrus{
		_base: base,
		_log:  log,
	}
}

// F implements TestInterface
func (_d TestInterfaceWithLogrus) F(ctx context.Context, a1 string, a2 ...string) (result1 string, result2 string, err error) {
	_d._log.WithFields(logrus.Fields(map[string]interface{}{
		"ctx": ctx,
		"a1":  a1,
		"a2":  a2})).Debug("TestInterfaceWithLogrus: calling F")
	defer func() {
		if err != nil {
			_d._log.WithFields(logrus.Fields(map[string]interface{}{
				"result1": result1,
				"result2": result2,
				"err":     err})).Error("TestInterfaceWithLogrus: method F returned an error")
		} else {
			_d._log.WithFields(logrus.Fields(map[string]interface{}{
				"result1": result1,
				"result2": result2,
				"err":     err})).Debug("TestInterfaceWithLogrus: method F finished")
		}
	}()
	return _d._base.F(ctx, a1, a2...)
}

// NoError implements TestInterface
func (_d TestInterfaceWithLogrus) NoError(s1 string) (s2 string) {
	_d._log.WithFields(logrus.Fields(map[string]interface{}{
		"s1": s1})).Debug("TestInterfaceWithLogrus: calling NoError")
	defer func() {
		_d._log.WithFields(logrus.Fields(map[string]interface{}{
			"s2": s2})).Debug("TestInterfaceWithLogrus: method NoError finished")
	}()
	return _d._base.NoError(s1)
}

// NoParamsOrResults implements TestInterface
func (_d TestInterfaceWithLogrus) NoParamsOrResults() {
	_d._log.Debug("TestInterfaceWithLogrus: calling NoParamsOrResults")
	defer func() {
		_d._log.Debug("TestInterfaceWithLogrus: NoParamsOrResults finished")
	}()
	_d._base.NoParamsOrResults()
	return
}
