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

func Test_WithConfig_json(t *testing.T) {
	o := options{log: logrus.New()}

	err := WithConfig(logger.Config{
		Verbosity: "error",
		Formatter: "json",
		WithColor: false,
	})(&o)
	require.NoError(t, err)

	assert.Equal(t, logrus.ErrorLevel, o.log.Level)
	assert.IsType(t, new(logrus.JSONFormatter), o.log.Formatter)
}

func Test_WithConfig_console(t *testing.T) {
	o := options{log: logrus.New()}

	err := WithConfig(logger.Config{
		Verbosity: "error",
		Formatter: "console",
		WithColor: false,
	})(&o)
	require.NoError(t, err)

	assert.Equal(t, logrus.ErrorLevel, o.log.Level)
	assert.IsType(t, new(logrus.TextFormatter), o.log.Formatter)
}

func Test_WithConfig_error(t *testing.T) {
	t.Run("unparsable level", func(t *testing.T) {
		o := options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Formatter: "json",
			Verbosity: "boum",
		})(&o)
		require.Error(t, err)
	})

	t.Run("unknown format", func(t *testing.T) {
		o := options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Verbosity: "error",
			Formatter: "boum",
		})(&o)
		require.Error(t, err)
	})

	t.Run("output stdout", func(t *testing.T) {
		o := options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Formatter: "json",
			Output:    "stdout",
		})(&o)
		require.NoError(t, err)
		assert.Equal(t, os.Stdout, o.log.Out)
	})

	t.Run("output stderr", func(t *testing.T) {
		o := options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Formatter: "json",
			Output:    "stderr",
		})(&o)
		require.NoError(t, err)
		assert.Equal(t, os.Stderr, o.log.Out)
	})

	t.Run("output empty", func(t *testing.T) {
		o := options{log: logrus.New()}

		err := WithConfig(logger.Config{
			Formatter: "json",
		})(&o)
		require.NoError(t, err)
		assert.Equal(t, os.Stderr, o.log.Out)
	})
}

func Test_WithLevel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		o := options{log: logrus.New()}
		err := WithLevel(logger.LevelError)(&o)
		require.NoError(t, err)
		assert.Equal(t, logrus.ErrorLevel, o.log.Level)
	})
	t.Run("fail", func(t *testing.T) {
		o := options{log: logrus.New()}
		err := WithLevel(logger.Level(42))(&o)
		require.Error(t, err)
	})
}

func Test_WithConsoleFormatter(t *testing.T) {
	o := options{log: logrus.New()}

	err := WithConsoleFormatter(true)(&o)
	require.NoError(t, err)
	assert.IsType(t, new(logrus.TextFormatter), o.log.Formatter)

	err = WithConsoleFormatter(false)(&o)
	require.NoError(t, err)
	assert.IsType(t, new(logrus.TextFormatter), o.log.Formatter)
}

func Test_WithJSONFormatter(t *testing.T) {
	o := options{log: logrus.New()}
	err := WithJSONFormatter()(&o)
	require.NoError(t, err)
	assert.IsType(t, new(logrus.JSONFormatter), o.log.Formatter)
}

func Test_WithOutput(t *testing.T) {
	var (
		o         = options{log: logrus.New()}
		_, writer = io.Pipe()
	)

	err := WithOutput(writer)(&o)

	require.NoError(t, err)
	assert.Equal(t, writer, o.log.Out)
}

func Test_WithInstance(t *testing.T) {
	var (
		o options
		l = logrus.New()
	)

	err := WithInstance(l)(&o)

	require.NoError(t, err)
	assert.Equal(t, l, o.log)
}

func Test_WithoutTime(t *testing.T) {
	t.Run("with text formatter", func(t *testing.T) {
		o := options{log: logrus.New()}
		require.NoError(t, WithConsoleFormatter(false)(&o))
		require.IsType(t, o.log.Formatter, &logrus.TextFormatter{})
		assert.False(t, (o.log.Formatter.(*logrus.TextFormatter)).DisableTimestamp)
		require.NoError(t, WithoutTime()(&o))
		assert.True(t, (o.log.Formatter.(*logrus.TextFormatter)).DisableTimestamp)
	})

	t.Run("with json formatter", func(t *testing.T) {
		o := options{log: logrus.New()}
		require.NoError(t, WithJSONFormatter()(&o))
		require.IsType(t, o.log.Formatter, &logrus.JSONFormatter{})
		assert.False(t, (o.log.Formatter.(*logrus.JSONFormatter)).DisableTimestamp)
		require.NoError(t, WithoutTime()(&o))
		assert.True(t, (o.log.Formatter.(*logrus.JSONFormatter)).DisableTimestamp)
	})

	t.Run("with unhandled formatter", func(t *testing.T) {
		o := options{log: logrus.New()}
		o.log.Formatter = nil
		require.Error(t, WithoutTime()(&o))
	})
}
