package logrus

import (
	"encoding/json"
	"fmt"
	stdlog "log"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/logger"
)

func newDeterministic() *Logrus {
	var log = New()

	log.log.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
		PrettyPrint:      false,
	}
	return log
}

func TestLogrusImplementLogger(t *testing.T) {
	var i interface{} = new(Logrus)
	if _, ok := i.(logger.Logger); !ok {
		t.Fatalf("expected %t to implement Logger", i)
	}
}

func TestConvertLevel(t *testing.T) {
	var tests = map[string]struct {
		expectedFailure     bool
		level               logger.Level
		expectedLogrusLevel logrus.Level
	}{
		"debug": {
			level:               logger.LevelDebug,
			expectedLogrusLevel: logrus.DebugLevel,
		}, "warn": {
			level:               logger.LevelWarn,
			expectedLogrusLevel: logrus.WarnLevel,
		}, "info": {
			level:               logger.LevelInfo,
			expectedLogrusLevel: logrus.InfoLevel,
		}, "error": {
			level:               logger.LevelError,
			expectedLogrusLevel: logrus.ErrorLevel,
		}, "failure": {
			level:           logger.Level(42),
			expectedFailure: true,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logrusLvl, err := convertLevel(test.level)
			if test.expectedFailure {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedLogrusLevel, logrusLvl)
			}
		})
	}
}

func TestRedirectStdLog(t *testing.T) {
	const imalog = "imalog"
	var expectedOutput = map[string]interface{}{
		"level":  logrus.ErrorLevel.String(),
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

	t.Run("using logrus", func(t *testing.T) {
		// redirect stdlog to zap
		outputRaw, err := logger.CaptureOutput(func() {
			var (
				log     = newDeterministic()
				restore = logger.RedirectStdLog(log, logger.LevelError)
			)
			defer restore()

			stdlog.Println("i'm a log")
		})
		require.NoError(t, err)

		var output map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(outputRaw), &output))
		assert.Equal(t, expectedOutput, output)
	})
}

func TestLogrus_SetLevel(t *testing.T) {
	var log = New()

	t.Run("nominal", func(t *testing.T) {
		err := log.SetLevel(logger.LevelWarn)
		require.NoError(t, err)
		assert.Equal(t, logrus.WarnLevel, log.log.Level)
	})

	t.Run("error", func(t *testing.T) {
		err := log.SetLevel(logger.Level(42))
		require.Error(t, err)
	})
}
