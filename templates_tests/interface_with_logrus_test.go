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

func TestTestInterfaceWithLogrus_F(t *testing.T) {
	t.Run("method returns an error", func(t *testing.T) {
		errUnexpected := errors.New("unexpected error")
		impl := &testImpl{err: errUnexpected, r1: "1", r2: "2"}

		buf := bytes.NewBuffer([]byte{})

		logger := logrus.Logger{
			Out:       buf,
			Formatter: &logrus.JSONFormatter{},
			Level:     logrus.DebugLevel,
		}

		entry := logrus.NewEntry(&logger)

		wrapped := NewTestInterfaceWithLogrus(impl, entry)

		r1, r2, err := wrapped.F(context.Background(), "a1")

		assert.Error(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)

		callingRecord := make(map[string]interface{})
		decoder := json.NewDecoder(buf)

		err = decoder.Decode(&callingRecord)

		require.NoError(t, err)
		assert.Equal(t, "a1", callingRecord["a1"])
		assert.Nil(t, callingRecord["a2"])
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
		}

		entry := logrus.NewEntry(&logger)

		wrapped := NewTestInterfaceWithLogrus(impl, entry)

		r1, r2, err := wrapped.F(context.Background(), "a1")

		assert.NoError(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)

		callingRecord := make(map[string]interface{})
		decoder := json.NewDecoder(buf)

		err = decoder.Decode(&callingRecord)

		require.NoError(t, err)
		assert.Equal(t, "a1", callingRecord["a1"])
		assert.Nil(t, callingRecord["a2"])
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
