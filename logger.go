// Package logger exposes a unique interface that hides the underlying implementation
// and exposes useful function to switch back to a writer, or to a standard logger.
package logger

// Logger defines the way logs can be handled.
type Logger interface {
	// Update apply the configuration on the logger.
	SetLevel(Level) error

	// Debug logs a message at the 'debug' level.
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	// Info logs a message at the 'info' level.
	Info(args ...interface{})
	Infof(format string, args ...interface{})

	// Info logs a message at the 'warn' level.
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	// Error logs a message at the 'error' level.
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	// WithField adds a field to the logging context.
	WithField(key string, value interface{}) Logger
	// WithFields adds multiple fields to the logging context.
	WithFields(fields map[string]interface{}) Logger
	// WithError adds an error field to the logging context.
	WithError(err error) Logger
}

// FieldErrorKey is the name of the field set by WithError.
const FieldErrorKey = "error"
