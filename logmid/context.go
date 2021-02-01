package logmid

import (
	"context"
	"fmt"
)

type customCtx string

const (
	ctxLogFieldsKey customCtx = "fields"
	ctxLogErrorsKey customCtx = "errors"
)

// AddFieldInContext add a log field to the context.
func AddFieldInContext(ctx context.Context, key string, value interface{}) {
	if fields, ok := ctx.Value(ctxLogFieldsKey).(map[string]interface{}); ok {
		fields[key] = value
	}
}

// AddErrorInContext add an error field to the context.
// If there is already an error in the context, the
// original error will be wrap with the reason (err.Error()).
func AddErrorInContext(ctx context.Context, err error) {
	if errCtx, ok := ctx.Value(ctxLogErrorsKey).(*error); ok {
		if *errCtx == nil {
			*errCtx = err
		} else {
			*errCtx = fmt.Errorf("%s: %w", err.Error(), *errCtx)
		}
	}
}
