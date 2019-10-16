// Package logrus implements the logger.Logger interface using sirupsen/logrus implementation.
package logrus

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/krostar/logger"
)

// Logrus implements Logger interface.
type Logrus struct {
	log *logrus.Logger
	logrus.FieldLogger
}

// New returns a new logrus instance.
func New(opts ...Option) (*Logrus, error) {
	var o options

	o.log = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, fmt.Errorf("unable to apply config: %w", err)
		}
	}

	return &Logrus{
		log:         o.log,
		FieldLogger: o.log,
	}, nil
}

func convertLevel(level logger.Level) (logrus.Level, error) {
	var logrusLevel logrus.Level
	switch level {
	case logger.LevelDebug:
		logrusLevel = logrus.DebugLevel
	case logger.LevelInfo:
		logrusLevel = logrus.InfoLevel
	case logger.LevelWarn:
		logrusLevel = logrus.WarnLevel
	case logger.LevelError:
		logrusLevel = logrus.ErrorLevel
	default:
		return logrusLevel, errors.New("level conversion to logrus level impossible")
	}
	return logrusLevel, nil
}

// SetLevel applies a new level to a logger instance.
func (l *Logrus) SetLevel(level logger.Level) error {
	lvl, err := convertLevel(level)
	if err != nil {
		return fmt.Errorf("unable to convert level: %w", err)
	}
	l.log.Level = lvl
	return nil
}

// WithField implements Logger.WithField for logrus's logger.
func (l *Logrus) WithField(key string, value interface{}) logger.Logger {
	return &Logrus{
		log:         l.log,
		FieldLogger: l.FieldLogger.WithField(key, value),
	}
}

// WithFields implements Logger.WithFields for logrus's logger.
func (l *Logrus) WithFields(fields map[string]interface{}) logger.Logger {
	return &Logrus{
		log:         l.log,
		FieldLogger: l.FieldLogger.WithFields(fields),
	}
}

// WithError implements Logger.WithError for logrus's logger.
func (l *Logrus) WithError(err error) logger.Logger {
	if err != nil {
		return l.WithField(logger.FieldErrorKey, err.Error())
	}
	return l
}
