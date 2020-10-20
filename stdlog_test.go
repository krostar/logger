package logger

import (
	stdlog "log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ToStdLog(t *testing.T) {
	var (
		log    = NewInMemory(LevelDebug)
		stdLog = StdLog(log, LevelWarn)
	)

	stdLog.Println("Hello world")

	require.Len(t, log.Entries, 1)
	assert.Len(t, log.Entries[0].Args, 1)
	assert.Equal(t, log.Entries[0].Args[0], "Hello world")
}

func Test_RedirectStdLog(t *testing.T) {
	const imalog = "i'm a log"

	var (
		log     = NewInMemory(LevelDebug)
		restore = RedirectStdLog(log, LevelError)
	)

	stdlog.Println(imalog)

	require.Len(t, log.Entries, 1)
	assert.Equal(t, LevelError, log.Entries[0].Level)
	assert.Equal(t, log.Entries[0].Fields, map[string]interface{}{
		"stdlog": "unhandled call to standard log package",
	})
	require.Len(t, log.Entries[0].Args, 1)
	assert.Equal(t, imalog, log.Entries[0].Args[0])

	restore()

	stdlog.Println(imalog)
	require.Len(t, log.Entries, 1)
}
