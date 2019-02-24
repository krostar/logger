package logmid

import (
	"net/http"

	"github.com/krostar/httpinfo"
	"github.com/krostar/logger"
)

type (
	// OnRequestFunc defines the signature of the function called on each requests.
	OnRequestFunc func(r *http.Request)

	// LogAtLevelFunc defines the signature of the function that gives the level of the log function.
	LogAtLevelFunc func(r *http.Request) logger.Level
)

// Options stores the middleware configuration options.
type Options struct {
	message    string
	onRequest  []OnRequestFunc
	logAtLevel LogAtLevelFunc
}

// Option defines a way to apply an option to the options.
type Option func(opts *Options)

// WithCallback adds a function called each time a request is logged.
func WithCallback(fcts ...OnRequestFunc) Option {
	return func(opts *Options) {
		opts.onRequest = append(opts.onRequest, fcts...)
	}
}

// WithLogLevelFunc adds a function to set logger's level based on the request.
func WithLogLevelFunc(fct LogAtLevelFunc) Option {
	return func(opts *Options) {
		opts.logAtLevel = fct
	}
}

// WithDefaultFields adds some fields to the request's log.
func WithDefaultFields() Option {
	return WithCallback(func(r *http.Request) {
		var ctx = r.Context()

		AddFieldInContext(ctx, "uri", r.URL.RequestURI())
		AddFieldInContext(ctx, "user-agent", r.Header.Get("User-Agent"))
		AddFieldInContext(ctx, "remote-addr", r.RemoteAddr)

		if httpinfo.IsUsed(r) {
			AddFieldInContext(ctx, "latency", httpinfo.ExecutionTime(r))
			AddFieldInContext(ctx, "status", httpinfo.Status(r))
			AddFieldInContext(ctx, "content-length", httpinfo.BytesWrote(r))
			AddFieldInContext(ctx, "route", httpinfo.RouteUsed(r))
		}
	})
}

// WithMessage sets a custom message on each http request logs.
func WithMessage(message string) Option {
	return func(opts *Options) {
		opts.message = message
	}
}
