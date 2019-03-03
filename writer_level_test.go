package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriterLevel(t *testing.T) {
	var l = NewInMemory(LevelDebug)

	writer := WriterLevel(l, LevelInfo)
	require.NotNil(t, writer)

	_, err := writer.Write([]byte("hello world"))
	require.NoError(t, err)

	require.Len(t, l.Entries, 1)
	assert.Equal(t, LevelInfo, l.Entries[0].Level)
	require.Len(t, l.Entries[0].Args, 1)
	require.Equal(t, "hello world", l.Entries[0].Args[0])
}
