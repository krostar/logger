// Package zap implements the logger.Logger interface using uber/zap implementation.
package zap

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/krostar/logger"
)

// Zap implements Logger interface.
type Zap struct {
	*zap.SugaredLogger
	level *zap.AtomicLevel
}

// New returns a new zap instance.
func New(opts ...Option) (logger.Logger, func() error, error) {
	config := config{
		Level: zapcore.InfoLevel,
		Zap: zap.Config{
			Development:       false,
			DisableCaller:     true,
			DisableStacktrace: true,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stderr"},
			Encoding:          "json",
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey: "msg",
				LineEnding: zapcore.DefaultLineEnding,

				LevelKey:    "lvl",
				EncodeLevel: zapcore.LowercaseLevelEncoder,

				TimeKey:        "time",
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,

				CallerKey:    "caller",
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
		},
	}

	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, nil, fmt.Errorf("unable to apply config: %w", err)
		}
	}

	atomiclevel := zap.NewAtomicLevelAt(config.Level)
	config.Zap.Level = atomiclevel

	logger, err := config.Zap.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create logger: %w", err)
	}

	return &Zap{
		level:         &atomiclevel,
		SugaredLogger: logger.Sugar(),
	}, logger.Sync, nil
}

func convertLevel(level logger.Level) (zapcore.Level, error) {
	var zapLevel zapcore.Level
	switch level {
	case logger.LevelDebug:
		zapLevel = zapcore.DebugLevel
	case logger.LevelInfo:
		zapLevel = zapcore.InfoLevel
	case logger.LevelWarn:
		zapLevel = zapcore.WarnLevel
	case logger.LevelError:
		zapLevel = zapcore.ErrorLevel
	default:
		return zapLevel, errors.New("level conversion to zap level impossible")
	}
	return zapLevel, nil
}

// SetLevel applies a new level to a logger instance.
func (l *Zap) SetLevel(level logger.Level) error {
	zapLevel, err := convertLevel(level)
	if err != nil {
		return fmt.Errorf("logger level parsing to zap level failed: %w", err)
	}
	(l.level).SetLevel(zapLevel)
	return nil
}

// WithField implements Logger.WithField for Zap logger.
func (l *Zap) WithField(key string, value interface{}) logger.Logger {
	return &Zap{
		level:         l.level,
		SugaredLogger: l.With(key, value),
	}
}

// WithFields implements Logger.WithFields for Zap logger.
func (l *Zap) WithFields(fields map[string]interface{}) logger.Logger {
	var f []interface{}
	for key, value := range fields {
		f = append(f, []interface{}{key, value}...)
	}

	return &Zap{
		level:         l.level,
		SugaredLogger: l.With(f...),
	}
}

// WithError implements Logger.WithError for Zap logger.
func (l *Zap) WithError(err error) logger.Logger {
	if err != nil {
		return l.WithField(logger.FieldErrorKey, err.Error())
	}
	return l
}
