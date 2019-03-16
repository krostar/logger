package zap

import (
	"github.com/pkg/errors"
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
	var config = config{
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
		opt(&config)
	}

	atomiclevel := zap.NewAtomicLevelAt(config.Level)
	config.Zap.Level = atomiclevel

	logger, err := config.Zap.Build()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create logger")
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

// Level returns the actual level of the zap's logger.
func (l *Zap) Level() logger.Level {
	var lvl logger.Level

	switch l.level.Level() {
	case zapcore.DebugLevel:
		lvl = logger.LevelDebug
	case zapcore.InfoLevel:
		lvl = logger.LevelInfo
	case zapcore.WarnLevel:
		lvl = logger.LevelWarn
	case zapcore.ErrorLevel:
		lvl = logger.LevelError
	}

	return lvl
}

// SetLevel applies a new level to a logger instance.
func (l *Zap) SetLevel(level logger.Level) error {
	zapLevel, err := convertLevel(level)
	if err != nil {
		return errors.Wrap(err, "logger level parsing to zap level failed")
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
