package zap

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		OutputErr: "yoloerr",
	})(&cfg)

	assert.Equal(t, zapcore.ErrorLevel, cfg.Level)
	assert.Equal(t, "json", cfg.Zap.Encoding)
	assert.Equal(t, []string{"yolo"}, cfg.Zap.OutputPaths)
	assert.Equal(t, []string{"yoloerr"}, cfg.Zap.ErrorOutputPaths)
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

func TestWithErrOutputPaths(t *testing.T) {
	var cfg config
	WithErrOutputPaths([]string{"yolo", "yili"})(&cfg)
	assert.Equal(t, []string{"yolo", "yili"}, cfg.Zap.ErrorOutputPaths)
}
