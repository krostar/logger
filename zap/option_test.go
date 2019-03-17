package zap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/krostar/logger"
)

func TestWithConfig(t *testing.T) {
	t.Run("success with json", func(t *testing.T) {
		var cfg config
		err := WithConfig(logger.Config{
			Verbosity: "error",
			Formatter: "json",
			WithColor: false,
			Output:    "yolo",
		})(&cfg)

		require.NoError(t, err)
		assert.Equal(t, zapcore.ErrorLevel, cfg.Level)
		assert.Equal(t, "json", cfg.Zap.Encoding)
		assert.Equal(t, []string{"yolo"}, cfg.Zap.OutputPaths)
	})

	t.Run("success with console", func(t *testing.T) {
		var cfg config
		err := WithConfig(logger.Config{
			Verbosity: "error",
			Formatter: "console",
			WithColor: false,
			Output:    "yolo",
		})(&cfg)

		require.NoError(t, err)
		assert.Equal(t, zapcore.ErrorLevel, cfg.Level)
		assert.Equal(t, "console", cfg.Zap.Encoding)
		assert.Equal(t, []string{"yolo"}, cfg.Zap.OutputPaths)
	})

	t.Run("unparsable level", func(t *testing.T) {
		var cfg config

		err := WithConfig(logger.Config{
			Verbosity: "boum",
		})(&cfg)
		require.Error(t, err)
	})

	t.Run("unknown formatter", func(t *testing.T) {
		var cfg config

		err := WithConfig(logger.Config{
			Verbosity: "error",
			Formatter: "boum",
		})(&cfg)
		require.Error(t, err)
	})
}

func TestWithLevel(t *testing.T) {
	var cfg config
	err := WithLevel(logger.LevelError)(&cfg)
	require.NoError(t, err)
	assert.Equal(t, zapcore.ErrorLevel, cfg.Level)
}

func TestWithConsoleFormatter(t *testing.T) {
	var cfg config

	err := WithConsoleFormatter(true)(&cfg)
	require.NoError(t, err)
	assert.Equal(t, "console", cfg.Zap.Encoding)
	err = WithConsoleFormatter(false)(&cfg)
	require.NoError(t, err)
	assert.Equal(t, "console", cfg.Zap.Encoding)
}

func TestWithJSONFormatter(t *testing.T) {
	var cfg config
	err := WithJSONFormatter()(&cfg)
	require.NoError(t, err)
	assert.Equal(t, "json", cfg.Zap.Encoding)
}

func TestWithOutputPaths(t *testing.T) {
	var cfg config
	err := WithOutputPaths([]string{"yolo", "yili"})(&cfg)
	require.NoError(t, err)
	assert.Equal(t, []string{"yolo", "yili"}, cfg.Zap.OutputPaths)
}

func TestWithZapConfig(t *testing.T) {
	var (
		cfg    config
		zapCfg = zap.Config{
			Development: true,
		}
	)
	err := WithZapConfig(zapCfg)(&cfg)
	require.NoError(t, err)
	assert.Equal(t, zapCfg, cfg.Zap)
}
