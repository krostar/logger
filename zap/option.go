package zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/krostar/logger"
)

type config struct {
	Level zapcore.Level
	Zap   zap.Config
}

// Option defines a function signature to update configuration.
type Option func(*config)

// WithConfig takes the logger configuration and applies it.
func WithConfig(cfg logger.Config) Option {
	var opts []Option

	// verbosity
	if lvl, err := logger.ParseLevel(cfg.Verbosity); err == nil {
		opts = append(opts, WithLevel(lvl))
	} else {
		fmt.Println("unable to apply verbosity")
	}

	// formatter
	switch cfg.Formatter {
	case "json":
		opts = append(opts, WithJSONFormatter())
	case "console":
		opts = append(opts, WithConsoleFormatter(cfg.WithColor))
	}

	// outputs
	if cfg.Output != "" {
		opts = append(opts, WithOutputPaths([]string{cfg.Output}))
	}
	if cfg.OutputErr != "" {
		opts = append(opts, WithErrOutputPaths([]string{cfg.OutputErr}))
	}

	return func(c *config) {
		for _, opt := range opts {
			opt(c)
		}
	}
}

// WithLevel configures the minimum level of the logger.
// It can late be updated with SetLevel.
func WithLevel(level logger.Level) Option {
	return func(c *config) {
		lvl, err := convertLevel(level)
		if err == nil {
			c.Level = lvl
		}
	}
}

// WithConsoleFormatter configures the format of the log output
//   to use "console" (cli) formatter.
func WithConsoleFormatter(colored bool) Option {
	return func(c *config) {
		c.Zap.Encoding = "console"
		if colored {
			c.Zap.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		} else {
			c.Zap.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		}
	}
}

// WithJSONFormatter configures the format of the log output
//   to use "json" formatter.
func WithJSONFormatter() Option {
	return func(c *config) {
		c.Zap.Encoding = "json"
		c.Zap.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
}

// WithOutputPaths configures the paths used to write logs to.
//   To use standart output, and error output, use stdout or stderr.
func WithOutputPaths(paths []string) Option {
	return func(c *config) {
		c.Zap.OutputPaths = paths
	}
}

// WithErrOutputPaths configures the paths used to write error logs to.
//   To use standart output, and error output, use stdout or stderr.
func WithErrOutputPaths(paths []string) Option {
	return func(c *config) {
		c.Zap.ErrorOutputPaths = paths
	}
}

// WithZapConfig applies zap configuration directly into the configuration.
func WithZapConfig(cfg zap.Config) Option {
	return func(c *config) {
		c.Zap = cfg
	}
}
