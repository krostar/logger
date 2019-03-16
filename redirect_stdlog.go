package logger

import (
	stdlog "log"
	"os"
)

// RedirectStdLog redirects standard logger calls to the underlying logger.
// This is heavily inspired by zap's way of doing the same thing.
func RedirectStdLog(l Logger, at Level) func() {
	var (
		oldFlags  = stdlog.Flags()
		oldPrefix = stdlog.Prefix()
	)
	stdlog.SetPrefix("")
	stdlog.SetFlags(0)

	stdlog.SetOutput(WriterLevel(
		l.WithField("stdlog", "unhandled call to standard log package"),
		at,
	))

	return func() {
		stdlog.SetFlags(oldFlags)
		stdlog.SetPrefix(oldPrefix)
		stdlog.SetOutput(os.Stderr)
	}
}
