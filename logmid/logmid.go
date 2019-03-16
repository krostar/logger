package logmid

import (
	"context"
	"net/http"

	"github.com/krostar/httpinfo"
	"github.com/krostar/logger"
)

// New returns a middleware that log requests.
func New(log logger.Logger, opts ...Option) func(http.Handler) http.Handler {
	var o = Options{
		message:    "http request",
		logAtLevel: defaultLogAtLevelFunc,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var (
				ctx    = r.Context()
				err    error
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
		case status >= 500:
			lvl = logger.LevelError
		case status >= 400 && status < 500:
			lvl = logger.LevelWarn
		default:
			lvl = logger.LevelInfo
		}
	} else {
		lvl = logger.LevelInfo
	}

	return lvl
}
