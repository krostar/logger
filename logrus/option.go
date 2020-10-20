package logrus

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/krostar/logger"
)

type options struct {
	log *logrus.Logger
}

// Option defines a function signature to update configuration.
type Option func(*options) error

// WithConfig takes the logger configuration and applies it.
func WithConfig(cfg logger.Config) Option {
	var opts []Option

	// verbosity
	if lvl, err := logger.ParseLevel(cfg.Verbosity); err == nil {
		opts = append(opts, WithLevel(lvl))
	} else {
		return func(c *options) error {
			return fmt.Errorf("unable to apply level %q, %w", cfg.Verbosity, err)
		}
	}

	// formatter
	switch cfg.Formatter {
	case "json":
		opts = append(opts, WithJSONFormatter())
	case "console":
		opts = append(opts, WithConsoleFormatter(cfg.WithColor))
	default:
		return func(c *options) error {
			return fmt.Errorf("unknown formatter %s", cfg.Formatter)
		}
	}

	// outputs
	opts = append(opts, withOutputStr(cfg.Output))

	// return all options
	return func(c *options) error {
		for _, opt := range opts {
			if err := opt(c); err != nil {
				return err
			}
		}
		return nil
	}
}

func withOutputStr(output string) Option {
	var opt = func(*options) error { return nil }

	if output == "" {
		return opt
	}

	switch output {
	case "stdout":
		opt = WithOutput(os.Stdout)
	case "stderr":
		opt = WithOutput(os.Stderr)
	default:
		f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return func(c *options) error {
				return fmt.Errorf("unable to open/create file %q: %w", output, err)
			}
		}
		opt = WithOutput(f)

		runtime.SetFinalizer(f, func(ff *os.File) {
			_ = ff.Sync()
			_ = ff.Close()
		})
	}

	return opt
}

// WithLevel configures the minimum level of the logger.
// It can later be updated with SetLevel.
func WithLevel(level logger.Level) Option {
	return func(o *options) error {
		lvl, err := convertLevel(level)
		if err != nil {
			return fmt.Errorf("failed to convert level: %w", err)
		}
		o.log.Level = lvl
		return nil
	}
}

// WithConsoleFormatter configures the format of the log output
// to use "console" (cli) formatter.
func WithConsoleFormatter(colored bool) Option {
	return func(o *options) error {
		var formatter logrus.TextFormatter
		if colored {
			formatter.DisableColors = false
		} else {
			formatter.DisableColors = true
		}
		o.log.Formatter = &formatter
		return nil
	}
}

// WithJSONFormatter configures the format of the log output
// to use "json" formatter.
func WithJSONFormatter() Option {
	return func(o *options) error {
		o.log.Formatter = new(logrus.JSONFormatter)
		return nil
	}
}

// WithOutput configures the writer used to write logs to.
func WithOutput(writer io.Writer) Option {
	return func(o *options) error {
		o.log.Out = writer
		return nil
	}
}

// WithInstance set logrus logger instance.
func WithInstance(log *logrus.Logger) Option {
	return func(o *options) error {
		o.log = log
		return nil
	}
}

// WithoutTime configures the logger to log without time.
// It only works with standard logrus formatter.
func WithoutTime() Option {
	return func(o *options) error {
		switch t := o.log.Formatter.(type) {
		case *logrus.TextFormatter:
			t.DisableTimestamp = true
		case *logrus.JSONFormatter:
			t.DisableTimestamp = true
		default:
			return fmt.Errorf("unhandled formatter %v", t)
		}
		return nil
	}
}
