package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryImplementLogger(t *testing.T) {
	var i interface{} = new(InMemory)
	if _, ok := i.(Logger); !ok {
		t.Fatalf("expected %t to implement Logger", i)
	}
}

func TestInMemory_SetLevel(t *testing.T) {
	var log = NewInMemory(LevelDebug)

	err := log.SetLevel(LevelError)
	assert.NoError(t, err)
	assert.Equal(t, LevelError, log.Level)
}

func TestInMemory_Log(t *testing.T) {
	var (
		log   = NewInMemory(LevelQuiet)
		tests = map[string]struct {
			logFunc       func(args ...interface{})
			logFFunc      func(format string, args ...interface{})
			minimalLevel  Level
			expectedLevel Level
		}{
			"debug": {
				logFunc:       log.Debug,
				logFFunc:      log.Debugf,
				minimalLevel:  LevelDebug,
				expectedLevel: LevelDebug,
			}, "info": {
				logFunc:       log.Info,
				logFFunc:      log.Infof,
				minimalLevel:  LevelInfo,
				expectedLevel: LevelInfo,
			}, "warn": {
				logFunc:       log.Warn,
				logFFunc:      log.Warnf,
				minimalLevel:  LevelWarn,
				expectedLevel: LevelWarn,
			}, "error": {
				logFunc:       log.Error,
				logFFunc:      log.Errorf,
				minimalLevel:  LevelError,
				expectedLevel: LevelError,
			},
		}
	)

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			// tests can't be ran in parallel as they share the same logger

			log.Reset(LevelQuiet)

			assert.Empty(t, log.Fields, "original logger should not contain any fields")
			assert.Empty(t, log.Entries, "original logger should not contain any entries")

			// try to log with not enough verbosity
			test.logFunc("anything", 42)
			test.logFFunc("anything %d", 42)
			assert.Empty(t, log.Fields, "no fields should have been set")
			assert.Empty(t, log.Entries, "the verbosity was not supposed to be high enough to display something")

			// now we should see logs
			log.SetLevel(test.minimalLevel) // nolint: errcheck, gosec

			// try again
			test.logFunc("another thing", 42)
			require.Len(t, log.Entries, 1)
			assert.Equal(t, test.expectedLevel, log.Entries[0].Level)
			assert.Empty(t, log.Entries[0].Format, "format should only be set in the logF variant")
			assert.Equal(t, log.Entries[0].Args, []interface{}{"another thing", 42})

			test.logFFunc("toto %d", 42)
			require.Len(t, log.Entries, 2)
			assert.Equal(t, test.expectedLevel, log.Entries[1].Level)
			assert.Equal(t, log.Entries[1].Format, "toto %d")
			assert.Equal(t, log.Entries[1].Args, []interface{}{42})
		})
	}
}

func TestInMemory_WithField(t *testing.T) {
	var log = NewInMemory(LevelDebug)

	assert.Empty(t, log.Entries)
	log.WithField("f", "ield").Debug("hello")
	assert.Len(t, log.Entries, 1)
	assert.Equal(t, map[string]interface{}{"f": "ield"}, log.Fields)
}

func TestInMemory_WithFields(t *testing.T) {
	var log = NewInMemory(LevelDebug)

	assert.Empty(t, log.Entries)
	log.WithFields(map[string]interface{}{"f": "ield"}).Debug("hello")
	assert.Len(t, log.Entries, 1)
	assert.Equal(t, map[string]interface{}{"f": "ield"}, log.Fields)
	assert.Equal(t, "ield", log.Fields["f"])
}

func TestInMemory_WithError(t *testing.T) {
	var log = NewInMemory(LevelDebug)

	assert.Empty(t, log.Entries)
	log.WithError(errors.New("eww")).Debug("hello")
	assert.Len(t, log.Entries, 1)
	assert.Equal(t, map[string]interface{}{FieldErrorKey: errors.New("eww")}, log.Fields)
}

func TestInMemory_Reset(t *testing.T) {
	var log = NewInMemory(LevelDebug)

	log.WithField("e", "f").Debug("debug")
	log.Reset(LevelInfo)

	assert.Empty(t, log.Entries)
	assert.Empty(t, log.Fields)
	assert.Equal(t, LevelInfo, log.Level)
}
