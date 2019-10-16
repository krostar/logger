package logmid

import (
	"context"
	"fmt"
)

type customCtx string

// nolint: gochecknoglobals
var (
	ctxLogFieldsKey = customCtx("fields")
	ctxLogErrorsKey = customCtx("errors")
)

// AddFieldInContext add a log field to the context.
func AddFieldInContext(ctx context.Context, key string, value interface{}) {
	var raw = ctx.Value(ctxLogFieldsKey)
	if fields, ok := raw.(map[string]interface{}); ok {
		fields[key] = value
	}
}

// AddErrorInContext add an error field to the context.
// If there is already an error in the context, the
//   original error will be wrap with the reason (err.Error()).
func AddErrorInContext(ctx context.Context, err error) {
	var raw = ctx.Value(ctxLogErrorsKey)
	if errCtx, ok := raw.(*error); ok {
		if *errCtx == nil {
			*errCtx = err
		} else {
			*errCtx = fmt.Errorf("%s: %w", err.Error(), *errCtx)
		}
	}
}
