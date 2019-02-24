package logger

import (
	"io/ioutil"
	"log"
	"os"
)

// CaptureOutput catpures and merge stdout and stderr.
// This function is voluntary not thread-safe as
// it belongs to the caller to make sure there should
// not be two logs that are written at the same time,
// because behavior in this case is undefined.
func CaptureOutput(writeFunc func()) string {
	var (
		stdout            = os.Stdout
		stderr            = os.Stderr
		reader, writer, _ = os.Pipe()
		output            = make(chan string)
	)

	// restore stdout and stderr to the original state
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()

	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)

	go func() {
		all, err := ioutil.ReadAll(reader) // nolint: errcheck
		if err != nil {
			close(output)
		}
		output <- string(all)
	}()

	writeFunc()
	writer.Close() // nolint: errcheck, gosec

	return <-output
}
