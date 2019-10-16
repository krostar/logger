package zap

import (
	"encoding/json"
	"errors"
	"fmt"
	stdlog "log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/krostar/logger"
)

func TestZapImplementLogger(t *testing.T) {
	var i interface{} = new(Zap)
	if _, ok := i.(logger.Logger); !ok {
		t.Fatalf("expected %t to implement Logger", i)
	}
}

func TestNew(t *testing.T) {
	log, f, err := New()
	assert.NotNil(t, log)
	assert.NotNil(t, f)
	assert.NoError(t, err)

	_, _, err = New(WithConfig(logger.Config{
		Formatter: "boum",
	}))
	require.Error(t, err)
}

func TestConvertLevel(t *testing.T) {
	var tests = map[string]struct {
		level            logger.Level
		expectedZapLevel zapcore.Level
		expectedFailure  bool
	}{
		"debug": {
			level:            logger.LevelDebug,
			expectedZapLevel: zapcore.DebugLevel,
		}, "warn": {
			level:            logger.LevelWarn,
			expectedZapLevel: zapcore.WarnLevel,
		}, "info": {
			level:            logger.LevelInfo,
			expectedZapLevel: zapcore.InfoLevel,
		}, "error": {
			level:            logger.LevelError,
			expectedZapLevel: zapcore.ErrorLevel,
		}, "failure": {
			level:           logger.Level(42),
			expectedFailure: true,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			zapLvl, err := convertLevel(test.level)
			if test.expectedFailure {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedZapLevel, zapLvl)
			}
		})
	}
}

func TestRedirectStdLog(t *testing.T) {
	const imalog = "imalog"
	var expectedOutput = map[string]interface{}{
		"level":  zapcore.ErrorLevel.String(),
		"msg":    "i'm a log",
		"stdlog": "unhandled call to standard log package",
	}

	t.Run("pilot", func(t *testing.T) {
		outputRaw, err := logger.CaptureOutput(func() {
			fmt.Print(imalog)
		})
		require.NoError(t, err)
		assert.Equal(t, imalog, outputRaw)
	})

	t.Run("using zap", func(t *testing.T) {
		// redirect stdlog to zap
		outputRaw, err := logger.CaptureOutput(func() {
			var log = &Zap{SugaredLogger: zap.NewExample().Sugar()}

			restore := logger.RedirectStdLog(log, logger.LevelError)
			defer restore()

			stdlog.Println("i'm a log")
		})
		require.NoError(t, err)

		var output map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(outputRaw), &output))
		assert.Equal(t, expectedOutput, output)
	})
}

func TestZap_SetLevel(t *testing.T) {
	var (
		err  error
		zLvl = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		log  = Zap{
			level: &zLvl,
		}
	)

	err = log.SetLevel(logger.LevelDebug)
	assert.NoError(t, err)
	assert.Equal(t, zapcore.DebugLevel, zLvl.Level())

	err = log.SetLevel(logger.Level(42))
	assert.Error(t, err)
	assert.Equal(t, zapcore.DebugLevel, zLvl.Level())
}

func TestZap_WithField(t *testing.T) {
	outputRaw, err := logger.CaptureOutput(func() {
		var log = &Zap{SugaredLogger: zap.NewExample().Sugar()}
		log.WithField("hello", "world").WithField("answer", 42).Warn("warn")
	})
	require.NoError(t, err)

	var output map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(outputRaw), &output))
	assert.Equal(t, map[string]interface{}{
		"level":  zapcore.WarnLevel.String(),
		"msg":    "warn",
		"hello":  "world",
		"answer": float64(42),
	}, output)
}

func TestZap_WithFields(t *testing.T) {
	outputRaw, err := logger.CaptureOutput(func() {
		var log = &Zap{SugaredLogger: zap.NewExample().Sugar()}
		log.
			WithFields(map[string]interface{}{"hello": "world"}).
			WithFields(map[string]interface{}{"answer": 42}).
			Warn("warn")
	})
	require.NoError(t, err)

	var output map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(outputRaw), &output))
	assert.Equal(t, map[string]interface{}{
		"level":  zapcore.WarnLevel.String(),
		"msg":    "warn",
		"hello":  "world",
		"answer": float64(42),
	}, output)
}

func TestZap_WithError(t *testing.T) {
	outputRaw, err := logger.CaptureOutput(func() {
		var log = &Zap{SugaredLogger: zap.NewExample().Sugar()}
		log.
			WithError(errors.New("eww1")).
			WithError(errors.New("eww2")).
			Warn("warn")
	})
	require.NoError(t, err)

	var output map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(outputRaw), &output))
	assert.Equal(t, map[string]interface{}{
		"level":              zapcore.WarnLevel.String(),
		"msg":                "warn",
		logger.FieldErrorKey: "eww2",
	}, output)
}
