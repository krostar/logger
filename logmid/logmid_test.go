package logmid

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/httpinfo"
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
	for _, entry := range log.Entries {
		for key := range entry.Fields {
			actualFieldsKeys = append(actualFieldsKeys, key)
		}
	}

	assert.ElementsMatch(t, expectedFieldsKeys, actualFieldsKeys)
	require.Len(t, log.Entries, 1)
	assert.Equal(t, expectedFields["key"], log.Entries[0].Fields["key"])
	assert.Equal(t, logger.LevelDebug, log.Entries[0].Level)
}

func TestDefaultLogAtLevelFunc(t *testing.T) {
	var (
		log   = logger.NewInMemory(logger.LevelDebug)
		tests = map[string]struct {
			status        int
			expectedLevel logger.Level
			useRecorder   bool
		}{
			"no recorder ok": {
				status:        http.StatusOK,
				expectedLevel: logger.LevelInfo,
				useRecorder:   false,
			},
			"no recorder nok": {
				status:        http.StatusInternalServerError,
				expectedLevel: logger.LevelInfo,
				useRecorder:   false,
			},
			"recorder >200<400": {
				status:        http.StatusOK,
				expectedLevel: logger.LevelInfo,
				useRecorder:   true,
			},
			"recorder >400<500": {
				status:        http.StatusBadRequest,
				expectedLevel: logger.LevelWarn,
				useRecorder:   true,
			},
			"recorder >500": {
				status:        http.StatusInternalServerError,
				expectedLevel: logger.LevelError,
				useRecorder:   true,
			},
		}
	)

	for name, test := range tests {
		var test = test

		t.Run(name, func(t *testing.T) {
			var (
				w                 = httptest.NewRecorder()
				r, _              = http.NewRequest("POST", "http://local/path?query", nil)
				h    http.Handler = http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
					rw.WriteHeader(test.status)
				})
			)

			h = New(log)(h)
			if test.useRecorder {
				h = httpinfo.Record()(h)
			}

			log.Reset()
			h.ServeHTTP(w, r)

			require.Len(t, log.Entries, 1)
			assert.Equal(t, test.expectedLevel, log.Entries[0].Level)
		})
	}
}
