package logger

import (
	"bytes"
	"io"
)

type wrappedWriter struct {
	logFunc func(...interface{})
}

func (ww *wrappedWriter) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	ww.logFunc(string(p))
	return len(p), nil
}

// WriterLevel returns a writer that writes as if
// LogAtLevelFunc(l) were called.
func WriterLevel(l Logger, at Level) io.Writer {
	return &wrappedWriter{logFunc: LogAtLevelFunc(l, at)}
}
