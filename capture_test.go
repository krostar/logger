package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureOutput(t *testing.T) {
	var stdout = os.Stdout
	var stderr = os.Stderr

	output := CaptureOutput(func() {
		fmt.Print("log")
	})
	assert.Equal(t, "log", output)

	assert.Equal(t, stdout, os.Stdout)
	assert.Equal(t, stderr, os.Stderr)
}
