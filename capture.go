package logger

import (
	"bytes"
	"fmt"
	"io"
	stdlog "log"
	"os"
)

// CaptureOutput catpures and merge stdout and stderr.
// This function is voluntary not thread-safe as it belongs
// to the caller to make sure there should not be two logs
// that are written from two different goroutines. Moreover
// os.Std{out,err} can't be modified from two different
// goroutines, because behavior in this case is undefined.
func CaptureOutput(writeFunc func()) (string, error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("unable to create pipe: %w", err)
	}

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() { // restore stdout and stderr to the original state
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		stdlog.SetOutput(os.Stderr)
	}()

	os.Stdout = writer
	os.Stderr = writer
	stdlog.SetOutput(writer)
	writeFunc()
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("unable to close writer: %w", err)
	}

	var output bytes.Buffer
	if _, err := io.Copy(&output, reader); err != nil {
		return "", fmt.Errorf("unable to read from reader: %w", err)
	}

	return output.String(), nil
}
