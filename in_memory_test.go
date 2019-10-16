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
func TestInMemory_Reset(t *testing.T) {
	var log = NewInMemory(LevelDebug)

	log.Info("hello world")
	assert.NotEmpty(t, log.Entries)

	log.Reset()
	assert.Empty(t, log.Entries)
}

func TestInMemory_SetLevel(t *testing.T) {
	var log = NewInMemory(LevelDebug)

	err := log.SetLevel(LevelError)
	assert.NoError(t, err)
	assert.Equal(t, LevelError, log.level)
}

func TestInMemory_WithField(t *testing.T) {
	var (
		lR  = NewInMemory(LevelDebug)
		llF = lR.WithField("hello", "world")
		lF  = llF.(*InMemory)
	)

	assert.NotEqual(t, lR, lF)
	assert.Equal(t, lR, lF.parent)
	assert.Equal(t, map[string]interface{}{
		"hello": "world",
	}, lF.fields)
}

func TestInMemory_WithFields(t *testing.T) {
	var (
		fields = map[string]interface{}{"hello": "world"}
		lR     = NewInMemory(LevelDebug)
		llF    = lR.WithFields(fields)
		lF     = llF.(*InMemory)
	)

	assert.NotEqual(t, lR, lF)
	assert.Equal(t, lR, lF.parent)
	assert.Equal(t, fields, lF.fields)
}

func TestInMemory_WithError(t *testing.T) {
	var (
		err = errors.New("hello world")
		lR  = NewInMemory(LevelDebug)
		llF = lR.WithError(err)
		lF  = llF.(*InMemory)
	)

	assert.NotEqual(t, lR, lF)
	assert.Equal(t, lR, lF.parent)
	assert.Equal(t, map[string]interface{}{
		FieldErrorKey: err,
	}, lF.fields)
}

func TestInMemory_log(t *testing.T) {
	var log = NewInMemory(LevelInfo)

	t.Run("root, inferior level", func(t *testing.T) {
		log.Reset()
		log.log(nil, LevelDebug, "%s", []interface{}{"hello"})

		assert.Empty(t, log.Entries)
	})

	t.Run("root, superior level", func(t *testing.T) {
		log.Reset()
		log.log(nil, LevelWarn, "%s", []interface{}{"hello"})

		require.Len(t, log.Entries, 1)
		assert.Equal(t, LevelWarn, log.Entries[0].Level)
		assert.Empty(t, log.Entries[0].Fields)
		assert.Equal(t, "%s", log.Entries[0].Format)
		assert.Equal(t, []interface{}{"hello"}, log.Entries[0].Args)
	})

	t.Run("root, with fields", func(t *testing.T) {
		log.Reset()
		log.log(map[string]interface{}{"hello": "world"}, LevelWarn, "%s", []interface{}{"toto"})

		require.Len(t, log.Entries, 1)
		assert.Equal(t, LevelWarn, log.Entries[0].Level)
		assert.Equal(t, map[string]interface{}{"hello": "world"}, log.Entries[0].Fields)
		assert.Equal(t, "%s", log.Entries[0].Format)
		assert.Equal(t, []interface{}{"toto"}, log.Entries[0].Args)
	})

	t.Run("child", func(t *testing.T) {
		log.Reset()

		child, _ := log.WithField("child", true).(*InMemory) // nolint: errcheck
		child.log(map[string]interface{}{"hello": "world"}, LevelWarn, "%s", []interface{}{"toto"})

		require.Len(t, log.Entries, 1)
		assert.Equal(t, LevelWarn, log.Entries[0].Level)
		assert.Equal(t, map[string]interface{}{
			"hello": "world",
			"child": true,
		}, log.Entries[0].Fields)
		assert.Equal(t, "%s", log.Entries[0].Format)
		assert.Equal(t, []interface{}{"toto"}, log.Entries[0].Args)
	})
}

func TestInMemory_Log(t *testing.T) {
	var (
		log   = NewInMemory(LevelQuiet)
		tests = map[string]struct {
			logFunc  func(args ...interface{})
			logFFunc func(format string, args ...interface{})
			level    Level
		}{
			"debug": {
				logFunc:  log.Debug,
				logFFunc: log.Debugf,
				level:    LevelDebug,
			}, "info": {
				logFunc:  log.Info,
				logFFunc: log.Infof,
				level:    LevelInfo,
			}, "warn": {
				logFunc:  log.Warn,
				logFFunc: log.Warnf,
				level:    LevelWarn,
			}, "error": {
				logFunc:  log.Error,
				logFFunc: log.Errorf,
				level:    LevelError,
			},
		}
	)

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			// tests can't be ran in parallel as they share the same logger
			log.Reset()
			log.SetLevel(LevelQuiet) // nolint: errcheck, gosec

			assert.Empty(t, log.fields, "original logger should not contain any fields")
			assert.Empty(t, log.Entries, "original logger should not contain any entries")

			// try to log with not enough verbosity
			test.logFunc("anything", 42)
			test.logFFunc("anything %d", 42)
			assert.Empty(t, log.Entries, "the verbosity was not supposed to be high enough to display something")

			// now we should see logs
			log.SetLevel(test.level) // nolint: errcheck, gosec

			// try again
			test.logFunc("another thing", 42)
			require.Len(t, log.Entries, 1)
			assert.Equal(t, test.level, log.Entries[0].Level)
			assert.Empty(t, log.Entries[0].Format, "format should only be set in the logF variant")
			assert.Equal(t, log.Entries[0].Args, []interface{}{"another thing", 42})

			test.logFFunc("toto %d", 42)
			require.Len(t, log.Entries, 2)
			assert.Equal(t, test.level, log.Entries[1].Level)
			assert.Equal(t, log.Entries[1].Format, "toto %d")
			assert.Equal(t, log.Entries[1].Args, []interface{}{42})
		})
	}
}
