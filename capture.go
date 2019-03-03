package logger

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
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
		return "", errors.Wrap(err, "unable to create pipe")
	}

	var (
		oldStdout = os.Stdout
		oldStderr = os.Stderr
	)
	defer func() { // restore stdout and stderr to the original state
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		log.SetOutput(os.Stderr)
	}()

	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)

	writeFunc()

	if err := writer.Close(); err != nil {
		return "", errors.Wrap(err, "unable to close writer")
	}

	var output bytes.Buffer
	if _, err := io.Copy(&output, reader); err != nil {
		return "", errors.Wrap(err, "unable to read from reader")
	}

	return output.String(), nil
}
