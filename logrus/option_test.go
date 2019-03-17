package logrus

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/logger"
)

func TestWithConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var o = options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Verbosity: "error",
			Formatter: "json",
			WithColor: false,
		})(&o)
		require.NoError(t, err)

		assert.Equal(t, logrus.ErrorLevel, o.log.Level)
		assert.IsType(t, new(logrus.JSONFormatter), o.log.Formatter)
	})

	t.Run("success with console", func(t *testing.T) {
		var o = options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Verbosity: "error",
			Formatter: "console",
			WithColor: false,
		})(&o)
		require.NoError(t, err)

		assert.Equal(t, logrus.ErrorLevel, o.log.Level)
		assert.IsType(t, new(logrus.TextFormatter), o.log.Formatter)
	})

	t.Run("unparsable level", func(t *testing.T) {
		var o = options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Verbosity: "boum",
		})(&o)
		require.Error(t, err)
	})

	t.Run("unknown format", func(t *testing.T) {
		var o = options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Verbosity: "error",
			Formatter: "boum",
		})(&o)
		require.Error(t, err)
	})
}

func TestWithLevel(t *testing.T) {
	var o = options{log: logrus.New()}
	err := WithLevel(logger.LevelError)(&o)
	require.NoError(t, err)
	assert.Equal(t, logrus.ErrorLevel, o.log.Level)
}

func TestWithConsoleFormatter(t *testing.T) {
	var o = options{log: logrus.New()}

	err := WithConsoleFormatter(true)(&o)
	require.NoError(t, err)
	assert.IsType(t, new(logrus.TextFormatter), o.log.Formatter)

	err = WithConsoleFormatter(false)(&o)
	require.NoError(t, err)
	assert.IsType(t, new(logrus.TextFormatter), o.log.Formatter)
}

func TestWithJSONFormatter(t *testing.T) {
	var o = options{log: logrus.New()}
	err := WithJSONFormatter()(&o)
	require.NoError(t, err)
	assert.IsType(t, new(logrus.JSONFormatter), o.log.Formatter)
}

func TestWithOutput(t *testing.T) {
	var (
		o         = options{log: logrus.New()}
		_, writer = io.Pipe()
	)

	err := WithOutput(writer)(&o)

	require.NoError(t, err)
	assert.Equal(t, writer, o.log.Out)
}

func TestWithInstance(t *testing.T) {
	var (
		o options
		l = logrus.New()
	)

	err := WithInstance(l)(&o)

	require.NoError(t, err)
	assert.Equal(t, l, o.log)
}
