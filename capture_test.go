package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CaptureOutput(t *testing.T) {
	stdout := os.Stdout
	stderr := os.Stderr

	output, err := CaptureOutput(func() {
		fmt.Print("log")
	})
	require.NoError(t, err)
	assert.Equal(t, "log", output)

	assert.Equal(t, stdout, os.Stdout)
	assert.Equal(t, stderr, os.Stderr)
}
