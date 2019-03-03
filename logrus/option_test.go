package logrus

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/krostar/logger"
)

func TestWithConfig(t *testing.T) {
}

func TestWithLevel(t *testing.T) {
	var o = options{log: logrus.New()}
	WithLevel(logger.LevelError)(&o)
	assert.Equal(t, logrus.ErrorLevel, o.log.Level)
}

func TestWithConsoleFormatter(t *testing.T) {
	var o = options{log: logrus.New()}
	WithConsoleFormatter(true)(&o)
	assert.IsType(t, new(logrus.TextFormatter), o.log.Formatter)
}

func TestWithJSONFormatter(t *testing.T) {
	var o = options{log: logrus.New()}
	WithJSONFormatter()(&o)
	assert.IsType(t, new(logrus.JSONFormatter), o.log.Formatter)
}

func TestWithOutput(t *testing.T) {
	var (
		o         = options{log: logrus.New()}
		_, writer = io.Pipe()
	)

	WithOutput(writer)(&o)

	assert.Equal(t, writer, o.log.Out)
}

func TestWithInstance(t *testing.T) {
	var (
		o options
		l = logrus.New()
	)

	WithInstance(l)(&o)

	assert.Equal(t, l, o.log)
}
