package templatestests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type contextKey struct{}

var key = contextKey{}

func TestTestInterfaceWithLogrus_F(t *testing.T) {
	t.Run("method returns an error", func(t *testing.T) {
		errUnexpected := errors.New("unexpected error")
		impl := &testImpl{err: errUnexpected, r1: "1", r2: "2"}

		buf := bytes.NewBuffer([]byte{})

		logger := logrus.Logger{
			Out:       buf,
			Formatter: &logrus.JSONFormatter{},
			Level:     logrus.DebugLevel,
			Hooks:     map[logrus.Level][]logrus.Hook{},
		}

		logger.AddHook(&ContextHook{})
		entry := logrus.NewEntry(&logger)

		wrapped := NewTestInterfaceWithLogrus(impl, entry)

		ctx := context.WithValue(context.Background(), key, "it does")
		r1, r2, err := wrapped.F(ctx, "a1")

		assert.Error(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)

		callingRecord := make(map[string]interface{})
		decoder := json.NewDecoder(buf)

		err = decoder.Decode(&callingRecord)

		require.NoError(t, err)
		assert.Equal(t, "a1", callingRecord["a1"])
		assert.Nil(t, callingRecord["a2"])
		assert.EqualValues(t, "it does", callingRecord["has_context"])
		assert.Equal(t, "TestInterfaceWithLogrus: calling F", callingRecord["msg"])

		finishedRecord := make(map[string]interface{})

		err = decoder.Decode(&finishedRecord)
		require.NoError(t, err)

		assert.Equal(t, "TestInterfaceWithLogrus: method F returned an error", finishedRecord["msg"])
		assert.Equal(t, "unexpected error", finishedRecord["err"])
		assert.Equal(t, "1", finishedRecord["result1"])
		assert.Equal(t, "2", finishedRecord["result2"])
	})

	t.Run("method finished successfully", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2"}

		buf := bytes.NewBuffer([]byte{})

		logger := logrus.Logger{
			Out:       buf,
			Formatter: &logrus.JSONFormatter{},
			Level:     logrus.DebugLevel,
			Hooks:     map[logrus.Level][]logrus.Hook{},
		}

		logger.AddHook(&ContextHook{})
		entry := logrus.NewEntry(&logger)

		wrapped := NewTestInterfaceWithLogrus(impl, entry)

		ctx := context.WithValue(context.Background(), key, "yes")
		r1, r2, err := wrapped.F(ctx, "a1")

		assert.NoError(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)

		callingRecord := make(map[string]interface{})
		decoder := json.NewDecoder(buf)

		err = decoder.Decode(&callingRecord)

		require.NoError(t, err)
		assert.Equal(t, "a1", callingRecord["a1"])
		assert.Nil(t, callingRecord["a2"])
		assert.EqualValues(t, "yes", callingRecord["has_context"])
		assert.Equal(t, "TestInterfaceWithLogrus: calling F", callingRecord["msg"])

		finishedRecord := make(map[string]interface{})

		err = decoder.Decode(&finishedRecord)
		require.NoError(t, err)

		assert.Equal(t, "TestInterfaceWithLogrus: method F finished", finishedRecord["msg"])
		assert.Equal(t, "1", finishedRecord["result1"])
		assert.Equal(t, "2", finishedRecord["result2"])
		assert.Nil(t, finishedRecord["err"])

		errorMessage, ok := finishedRecord["err"]
		assert.True(t, ok)
		assert.Nil(t, errorMessage)
	})
}

type ContextHook struct {
}

func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		entry.Data["has_context"] = entry.Context.Value(key)
	}
	return nil
}

func (hook *ContextHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
