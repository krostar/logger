package logger

import (
	"log"
	"os"
)

// RedirectStdLog redirects standard logger calls to the underlying logger.
// This is heavily inspired by zap's way of doing the same thing.
func RedirectStdLog(l Logger, at Level) func() {
	var (
		oldFlags  = log.Flags()
		oldPrefix = log.Prefix()
	)
	log.SetPrefix("")
	log.SetFlags(0)

	log.SetOutput(WriterLevel(
		l.WithField("stdlog", "unhandled call to standard log package"),
		at,
	))

	return func() {
		log.SetFlags(oldFlags)
		log.SetPrefix(oldPrefix)
		log.SetOutput(os.Stderr)
	}
}
