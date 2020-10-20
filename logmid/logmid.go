// Package logmid is a net/http middleware that uses httpinfo and
// logger.Logger to log HTTP requests.
package logmid

import (
	"context"
	"net/http"

	"github.com/krostar/httpinfo"

	"github.com/krostar/logger"
)

// New returns a middleware that log requests.
func New(log logger.Logger, opts ...Option) func(http.Handler) http.Handler {
	o := Options{
		message:    "http request",
		logAtLevel: defaultLogAtLevelFunc,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var (
				err    error
				ctx    = r.Context()
				fields = make(map[string]interface{})
			)

			ctx = context.WithValue(ctx, ctxLogErrorsKey, &err)
			ctx = context.WithValue(ctx, ctxLogFieldsKey, fields)
			r = r.WithContext(ctx)

			next.ServeHTTP(rw, r)

			for _, fct := range o.onRequest {
				fct(r)
			}

			requestLogger := log.WithFields(fields)
			if err != nil {
				requestLogger = requestLogger.WithError(err)
			}

			logger.LogAtLevelFunc(requestLogger, o.logAtLevel(r))(o.message)
		})
	}
}

func defaultLogAtLevelFunc(r *http.Request) logger.Level {
	var lvl logger.Level

	if httpinfo.IsUsed(r) {
		switch status := httpinfo.Status(r); {
		case status >= http.StatusInternalServerError:
			lvl = logger.LevelError
		case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
			lvl = logger.LevelWarn
		default:
			lvl = logger.LevelInfo
		}
	} else {
		lvl = logger.LevelInfo
	}

	return lvl
}
