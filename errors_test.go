package gowrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandLineError_Error(t *testing.T) {
	assert.Equal(t, "error", CommandLineError("error").Error())
}
