package logger

import (
	"errors"
	"fmt"
	"strings"
)

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// LevelDebug logs are typically voluminous, and are usually disabled in production.
	LevelDebug Level = iota
	// LevelInfo is the default logging priority.
	LevelInfo
	// LevelWarn logs are more important than Info, and may need individual human review.
	LevelWarn
	// LevelError logs are high-priority and should require a human review.
	LevelError
	// LevelQuiet hide everything.
	LevelQuiet

	levelDebugStr = "debug"
	levelInfoStr  = "info"
	levelWarnStr  = "warn"
	levelErrorStr = "error"
	levelQuietStr = "quiet"
)

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return levelDebugStr
	case LevelInfo:
		return levelInfoStr
	case LevelWarn:
		return levelWarnStr
	case LevelError:
		return levelErrorStr
	case LevelQuiet:
		return levelQuietStr
	default:
		return fmt.Sprintf("unknown level (%d)", l)
	}
}

// ParseLevel converts a string representation of a level to a level type.
func ParseLevel(levelStr string) (Level, error) {
	var l Level

	levelStr = strings.ToLower(levelStr)

	switch levelStr {
	case levelDebugStr:
		l = LevelDebug
	case levelInfoStr, "": // make the zero value useful
		l = LevelInfo
	case levelWarnStr:
		l = LevelWarn
	case levelErrorStr:
		l = LevelError
	case levelQuietStr:
		l = LevelQuiet
	default:
		return l, errors.New("unknown level")
	}

	return l, nil
}

// LogAtLevelFunc returns a function that can log to the provided level.
func LogAtLevelFunc(log Logger, l Level) func(...interface{}) {
	switch l {
	case LevelDebug:
		return log.Debug
	case LevelInfo:
		return log.Info
	case LevelWarn:
		return log.Warn
	case LevelError:
		return log.Error
	case LevelQuiet:
		return func(...interface{}) {}
	default:
		return log.Info
	}
}

// LogFAtLevelFunc is the same as LogAtLevelFunc but with log formatting.
func LogFAtLevelFunc(log Logger, l Level) func(string, ...interface{}) {
	switch l {
	case LevelDebug:
		return log.Debugf
	case LevelInfo:
		return log.Infof
	case LevelWarn:
		return log.Warnf
	case LevelError:
		return log.Errorf
	case LevelQuiet:
		return func(string, ...interface{}) {}
	default:
		return log.Infof
	}
}
