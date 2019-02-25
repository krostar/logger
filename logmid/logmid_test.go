package logmid

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/krostar/logger"
)

func TestMiddleware(t *testing.T) {
	var (
		log            = logger.NewInMemory(logger.LevelDebug)
		w              = httptest.NewRecorder()
		r, _           = http.NewRequest("POST", "http://local/path?query", nil)
		expectedFields = map[string]interface{}{
			"key":                "value",
			logger.FieldErrorKey: "", // value will not be checked, only the key
		}
		handler = func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("wrote"))
			assert.NoError(t, err)
		}
	)
	r.Header.Set("User-Agent", "super-agent-test")
	r.RemoteAddr = "10.11.12.13"

	New(log,
		WithCallback(func(r *http.Request) {
			var ctx = r.Context()
			AddFieldInContext(ctx, "key", "value")
			AddErrorInContext(ctx, errors.New("eww"))
		}),
		WithLogLevelFunc(func(r *http.Request) logger.Level {
			return logger.LevelDebug
		}),
	)(http.HandlerFunc(handler)).ServeHTTP(w, r)

	var expectedFieldsKeys []string
	for key := range expectedFields {
		expectedFieldsKeys = append(expectedFieldsKeys, key)
	}

	var actualFieldsKeys []string
	for key := range log.Fields {
		actualFieldsKeys = append(actualFieldsKeys, key)
	}

	assert.ElementsMatch(t, expectedFieldsKeys, actualFieldsKeys)
	assert.Equal(t, expectedFields["key"], log.Fields["key"])
	assert.Len(t, log.Entries, 1)
	assert.Equal(t, logger.LevelDebug, log.Entries[0].Level)
}
