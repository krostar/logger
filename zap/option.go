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
type Option func(*config) error

// WithConfig takes the logger configuration and applies it.
func WithConfig(cfg logger.Config) Option {
	var opts []Option

	// verbosity
	if lvl, err := logger.ParseLevel(cfg.Verbosity); err == nil {
		opts = append(opts, WithLevel(lvl))
	} else {
		return func(c *config) error {
			return fmt.Errorf("unable to apply level %q: %w", cfg.Verbosity, err)
		}
	}

	// formatter
	switch cfg.Formatter {
	case "json":
		opts = append(opts, WithJSONFormatter())
	case "console":
		opts = append(opts, WithConsoleFormatter(cfg.WithColor))
	default:
		return func(c *config) error {
			return fmt.Errorf("unknown formatter %s", cfg.Formatter)
		}
	}

	// outputs
	if cfg.Output != "" {
		opts = append(opts, WithOutputPaths([]string{cfg.Output}))
	}

	return func(c *config) error {
		for _, opt := range opts {
			if err := opt(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithLevel configures the minimum level of the logger.
// It can late be updated with SetLevel.
func WithLevel(level logger.Level) Option {
	return func(c *config) error {
		lvl, err := convertLevel(level)
		if err != nil {
			return err
		}
		c.Level = lvl
		return nil
	}
}

// WithConsoleFormatter configures the format of the log output
//   to use "console" (cli) formatter.
func WithConsoleFormatter(colored bool) Option {
	return func(c *config) error {
		c.Zap.Encoding = "console"
		if colored {
			c.Zap.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		} else {
			c.Zap.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		}
		return nil
	}
}

// WithJSONFormatter configures the format of the log output
//   to use "json" formatter.
func WithJSONFormatter() Option {
	return func(c *config) error {
		c.Zap.Encoding = "json"
		c.Zap.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		return nil
	}
}

// WithOutputPaths configures the paths used to write logs to.
//   To use standart output, and error output, use stdout or stderr.
func WithOutputPaths(paths []string) Option {
	return func(c *config) error {
		c.Zap.OutputPaths = paths
		return nil
	}
}

// WithZapConfig applies zap configuration directly into the configuration.
func WithZapConfig(cfg zap.Config) Option {
	return func(c *config) error {
		c.Zap = cfg
		return nil
	}
}
