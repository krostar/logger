package logger

import (
	stdlog "log"
	"os"
)

// StdLog returns a standard logger which will use
// the provided logger internally.
func StdLog(logger Logger, at Level) *stdlog.Logger {
	return stdlog.New(WriterLevel(logger, at), "", 0)
}

// RedirectStdLog redirects standard logger calls to the underlying logger.
// This is heavily inspired by zap's way of doing the same thing.
func RedirectStdLog(l Logger, at Level) func() {
	oldFlags := stdlog.Flags()
	oldPrefix := stdlog.Prefix()

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
