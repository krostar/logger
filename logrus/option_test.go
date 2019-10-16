package logrus

import (
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/logger"
)

func TestWithConfig_json(t *testing.T) {
	var o = options{log: logrus.New()}

	err := WithConfig(logger.Config{
		Verbosity: "error",
		Formatter: "json",
		WithColor: false,
	})(&o)
	require.NoError(t, err)

	assert.Equal(t, logrus.ErrorLevel, o.log.Level)
	assert.IsType(t, new(logrus.JSONFormatter), o.log.Formatter)
}

func TestWithConfig_console(t *testing.T) {
	var o = options{log: logrus.New()}

	err := WithConfig(logger.Config{
		Verbosity: "error",
		Formatter: "console",
		WithColor: false,
	})(&o)
	require.NoError(t, err)

	assert.Equal(t, logrus.ErrorLevel, o.log.Level)
	assert.IsType(t, new(logrus.TextFormatter), o.log.Formatter)
}

func TestWithConfig_error(t *testing.T) {
	t.Run("unparsable level", func(t *testing.T) {
		var o = options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Formatter: "json",
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

	t.Run("output stdout", func(t *testing.T) {
		var o = options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Formatter: "json",
			Output:    "stdout",
		})(&o)
		require.NoError(t, err)
		assert.Equal(t, os.Stdout, o.log.Out)
	})

	t.Run("output stderr", func(t *testing.T) {
		var o = options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Formatter: "json",
			Output:    "stderr",
		})(&o)
		require.NoError(t, err)
		assert.Equal(t, os.Stderr, o.log.Out)
	})

	t.Run("output empty", func(t *testing.T) {
		var o = options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Formatter: "json",
		})(&o)
		require.NoError(t, err)
		assert.Equal(t, os.Stderr, o.log.Out)
	})
}

func TestWithLevel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var o = options{log: logrus.New()}
		err := WithLevel(logger.LevelError)(&o)
		require.NoError(t, err)
		assert.Equal(t, logrus.ErrorLevel, o.log.Level)
	})
	t.Run("fail", func(t *testing.T) {
		var o = options{log: logrus.New()}
		err := WithLevel(logger.Level(42))(&o)
		require.Error(t, err)
	})
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
