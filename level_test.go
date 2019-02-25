package logger

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	var tests = map[string]struct {
		levelStr        string
		expectedLevel   Level
		expectedFailure bool
	}{
		"unknown level": {
			levelStr:        "lwfi",
			expectedFailure: true,
		}, "debug level": {
			levelStr:      levelDebugStr,
			expectedLevel: LevelDebug,
		}, "info level": {
			levelStr:      levelInfoStr,
			expectedLevel: LevelInfo,
		}, "warn level": {
			levelStr:      levelWarnStr,
			expectedLevel: LevelWarn,
		}, "error level": {
			levelStr:      levelErrorStr,
			expectedLevel: LevelError,
		}, "quiet level": {
			levelStr:      levelQuietStr,
			expectedLevel: LevelQuiet,
		}, "level with uppercase": {
			levelStr:      strings.ToUpper(levelInfoStr),
			expectedLevel: LevelInfo,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lvl, err := ParseLevel(test.levelStr)
			if test.expectedFailure {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expectedLevel, lvl)
		})
	}
}

func TestString(t *testing.T) {
	var tests = map[string]struct {
		level            Level
		expectedLevelStr string
	}{
		"unknown level": {
			level:            Level(9),
			expectedLevelStr: "unknown level (9)",
		}, "debug level": {
			level:            LevelDebug,
			expectedLevelStr: levelDebugStr,
		}, "info level": {
			level:            LevelInfo,
			expectedLevelStr: levelInfoStr,
		}, "warn level": {
			level:            LevelWarn,
			expectedLevelStr: levelWarnStr,
		}, "error level": {
			level:            LevelError,
			expectedLevelStr: levelErrorStr,
		}, "quiet level": {
			level:            LevelQuiet,
			expectedLevelStr: levelQuietStr,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedLevelStr, test.level.String())
		})
	}
}

func TestLogAtLevelFunc(t *testing.T) {
	var (
		log   = NewInMemory(LevelDebug)
		tests = map[string]struct {
			level             Level
			expectedNoEntries bool
		}{
			"debug level": {
				level: LevelDebug,
			}, "info level": {
				level: LevelInfo,
			}, "warn level": {
				level: LevelWarn,
			}, "error level": {
				level: LevelError,
			}, "quiet level": {
				level:             LevelQuiet,
				expectedNoEntries: true,
			},
		}
	)
	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			log.Reset(LevelDebug)

			LogAtLevelFunc(log, test.level)("log")

			if test.expectedNoEntries {
				require.Empty(t, log.Entries)
			} else {
				require.Len(t, log.Entries, 1)
				assert.Equal(t, test.level, log.Entries[0].Level)
				assert.Len(t, log.Entries[0].Args, 1)
				assert.Equal(t, "log", log.Entries[0].Args[0])
			}
		})
	}
}

func TestLogFAtLevelFunc(t *testing.T) {
	var (
		log   = NewInMemory(LevelDebug)
		tests = map[string]struct {
			level             Level
			expectedNoEntries bool
		}{
			"debug level": {
				level: LevelDebug,
			}, "info level": {
				level: LevelInfo,
			}, "warn level": {
				level: LevelWarn,
			}, "error level": {
				level: LevelError,
			}, "quiet level": {
				level:             LevelQuiet,
				expectedNoEntries: true,
			},
		}
	)
	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			log.Reset(LevelDebug)

			LogFAtLevelFunc(log, test.level)("log %d", 42)

			if test.expectedNoEntries {
				require.Empty(t, log.Entries)
			} else {
				require.Len(t, log.Entries, 1)
				assert.Equal(t, test.level, log.Entries[0].Level)
				assert.Equal(t, "log %d", log.Entries[0].Format)
				assert.Len(t, log.Entries[0].Args, 1)
				assert.Equal(t, 42, log.Entries[0].Args[0])
			}
		})
	}
}
