package zap

import (
	"encoding/json"
	"fmt"
	"log"
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
	z, f, err := New()
	assert.NotNil(t, z)
	assert.NotNil(t, f)
	assert.NoError(t, err)
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
			var z = &Zap{SugaredLogger: zap.NewExample().Sugar()}

			restore := logger.RedirectStdLog(z, logger.LevelError)
			defer restore()

			log.Println("i'm a log")
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
		z    = Zap{
			level: &zLvl,
		}
	)

	err = z.SetLevel(logger.LevelDebug)
	assert.NoError(t, err)
	assert.Equal(t, zapcore.DebugLevel, zLvl.Level())

	err = z.SetLevel(logger.Level(42))
	assert.Error(t, err)
	assert.Equal(t, zapcore.DebugLevel, zLvl.Level())
}
