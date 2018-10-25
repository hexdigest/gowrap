package templatestests

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestInterfaceWithLog_F(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2"}

		errLog := bytes.NewBuffer([]byte{})
		stdLog := bytes.NewBuffer([]byte{})

		wrapped := NewTestInterfaceWithLog(impl, stdLog, errLog)
		r1, r2, err := wrapped.F(context.Background(), "p1", "p2", "p3")
		require.NoError(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)

		assert.Empty(t, errLog.Bytes())
		assert.Contains(t, stdLog.String(), "TestInterfaceWithLog: calling F with params: context.Background p1 [p2 p3]")
		assert.Contains(t, stdLog.String(), "TestInterfaceWithLog: F returned results: 1 2 <nil>")
	})

	t.Run("error", func(t *testing.T) {
		impl := &testImpl{r1: "1", r2: "2", err: errors.New("failure")}

		errLog := bytes.NewBuffer([]byte{})
		stdLog := bytes.NewBuffer([]byte{})

		wrapped := NewTestInterfaceWithLog(impl, stdLog, errLog)
		r1, r2, err := wrapped.F(context.Background(), "p1", "p2", "p3")
		require.Error(t, err)
		assert.Equal(t, "1", r1)
		assert.Equal(t, "2", r2)

		assert.Contains(t, stdLog.String(), "TestInterfaceWithLog: calling F with params: context.Background p1 [p2 p3]")
		assert.Contains(t, errLog.String(), "")
	})
}
