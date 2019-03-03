package logger

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectStdLog(t *testing.T) {
	const imalog = "i'm a log"

	var (
		l       = NewInMemory(LevelDebug)
		restore = RedirectStdLog(l, LevelError)
	)

	log.Println(imalog)

	require.Len(t, l.Entries, 1)
	assert.Equal(t, LevelError, l.Entries[0].Level)
	assert.Equal(t, l.Entries[0].Fields, map[string]interface{}{
		"stdlog": "unhandled call to standard log package",
	})
	require.Len(t, l.Entries[0].Args, 1)
	assert.Equal(t, imalog, l.Entries[0].Args[0])

	restore()

	log.Println(imalog)
	require.Len(t, l.Entries, 1)
}
