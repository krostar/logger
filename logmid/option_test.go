package logmid

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krostar/httpinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/logger"
)

func Test_WithCallback(t *testing.T) {
	var opts Options
	WithCallback(func(r *http.Request) {})(&opts)
	assert.Len(t, opts.onRequest, 1)
}

func Test_WithLogLevelFunc(t *testing.T) {
	var opts Options
	WithLogLevelFunc(func(r *http.Request) logger.Level { return logger.LevelDebug })(&opts)
	assert.NotNil(t, opts.logAtLevel)
}

func Test_WithDefaultFields(t *testing.T) {
	log := logger.NewInMemory(logger.LevelDebug)
	mid := New(log, WithDefaultFields())

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://local/path?query", nil)

	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("wrote"))
		assert.NoError(t, err)
	}

	r.Header.Set("User-Agent", "homemade")
	r.RemoteAddr = "here"

	httpinfo.Record()(mid(http.HandlerFunc(handler))).ServeHTTP(w, r)

	require.Len(t, log.Entries, 1)
	entry := log.Entries[0]
	assert.Equal(t, "/path?query", entry.Fields["uri"])
	assert.Equal(t, "homemade", entry.Fields["user-agent"])
	assert.Equal(t, "here", entry.Fields["remote-addr"])
	// added thanks to httpinfo
	assert.Contains(t, entry.Fields, "latency")
	assert.Contains(t, entry.Fields, "status")
	assert.Contains(t, entry.Fields, "content-length")
	assert.Contains(t, entry.Fields, "route")
}

func Test_WithMessage(t *testing.T) {
	var opts Options
	WithMessage("hello world")(&opts)
	assert.Equal(t, "hello world", opts.message)
}
