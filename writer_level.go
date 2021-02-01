package logger

import (
	"bytes"
	"io"
)

type wrappedWriter struct {
	logFunc   func(...interface{})
	closeFunc func() error
}

func (ww *wrappedWriter) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	ww.logFunc(string(p))
	return len(p), nil
}

func (ww *wrappedWriter) Close() error {
	if ww.closeFunc != nil {
		return ww.closeFunc()
	}
	return nil
}

// WriterLevel returns a writer that writes as if
// LogAtLevelFunc(l) were called.
func WriterLevel(l Logger, at Level) io.Writer {
	return &wrappedWriter{logFunc: LogAtLevelFunc(l, at)}
}

// WriteCloserLevel returns a writer that writes as if
// LogAtLevelFunc(l) were called, with a closer.
func WriteCloserLevel(l Logger, close func() error, at Level) io.WriteCloser {
	return &wrappedWriter{logFunc: LogAtLevelFunc(l, at), closeFunc: close}
}
