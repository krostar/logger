package zap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/krostar/logger"
)

func TestWithConfig(t *testing.T) {
	var cfg config
	WithConfig(logger.Config{
		Verbosity: "error",
		Formatter: "json",
		WithColor: false,
		Output:    "yolo",
	})(&cfg)

	assert.Equal(t, zapcore.ErrorLevel, cfg.Level)
	assert.Equal(t, "json", cfg.Zap.Encoding)
	assert.Equal(t, []string{"yolo"}, cfg.Zap.OutputPaths)
}

func TestWithLevel(t *testing.T) {
	var cfg config
	WithLevel(logger.LevelError)(&cfg)
	assert.Equal(t, zapcore.ErrorLevel, cfg.Level)
}

func TestWithConsoleFormatter(t *testing.T) {
	var cfg config
	WithConsoleFormatter(true)(&cfg)
	assert.Equal(t, "console", cfg.Zap.Encoding)
}

func TestWithJSONFormatter(t *testing.T) {
	var cfg config
	WithJSONFormatter()(&cfg)
	assert.Equal(t, "json", cfg.Zap.Encoding)
}

func TestWithOutputPaths(t *testing.T) {
	var cfg config
	WithOutputPaths([]string{"yolo", "yili"})(&cfg)
	assert.Equal(t, []string{"yolo", "yili"}, cfg.Zap.OutputPaths)
}

func TestWithZapConfig(t *testing.T) {
	var (
		cfg    config
		zapCfg = zap.Config{
			Development: true,
		}
	)

	WithZapConfig(zapCfg)(&cfg)

	assert.Equal(t, zapCfg, cfg.Zap)
}
