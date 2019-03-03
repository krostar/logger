package logrus

import (
	"io"
	"os"
	"runtime"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/krostar/logger"
)

type options struct {
	log *logrus.Logger
}

// Option defines a function signature to update configuration.
type Option func(*options)

// WithConfig takes the logger configuration and applies it.
func WithConfig(cfg logger.Config) Option {
	var opts []Option

	// verbosity
	if lvl, err := logger.ParseLevel(cfg.Verbosity); err == nil {
		opts = append(opts, WithLevel(lvl))
	} else {
		panic(errors.Wrapf(err, "unable to parse level %q", cfg.Verbosity))
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
		f, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			panic(errors.Wrapf(err, "unable to open/create file %q", cfg.Output))
		}
		opts = append(opts, WithOutput(f))

		runtime.SetFinalizer(f, func(ff *os.File) {
			ff.Sync()  // nolint: errcheck, gosec
			ff.Close() // nolint: errcheck, gosec
		})
	}

	return func(c *options) {
		for _, opt := range opts {
			opt(c)
		}
	}
}

// WithLevel configures the minimum level of the logger.
// It can later be updated with SetLevel.
func WithLevel(level logger.Level) Option {
	return func(o *options) {
		lvl, err := convertLevel(level)
		if err == nil {
			o.log.Level = lvl
		}
	}
}

// WithConsoleFormatter configures the format of the log output
// to use "console" (cli) formatter.
func WithConsoleFormatter(colored bool) Option {
	return func(o *options) {
		var formatter logrus.TextFormatter
		if colored {
			formatter.DisableColors = false
		} else {
			formatter.DisableColors = true
		}
		o.log.Formatter = &formatter
	}
}

// WithJSONFormatter configures the format of the log output
// to use "json" formatter.
func WithJSONFormatter() Option {
	return func(o *options) {
		o.log.Formatter = new(logrus.JSONFormatter)
	}
}

// WithOutput configures the writer used to write logs to.
func WithOutput(writer io.Writer) Option {
	return func(o *options) {
		o.log.Out = writer
	}
}

// WithInstance set logrus logger instance.
func WithInstance(log *logrus.Logger) Option {
	return func(o *options) {
		o.log = log
	}
}
