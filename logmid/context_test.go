package logmid

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddFieldInContext(t *testing.T) {
	var (
		fields = make(map[string]interface{})
		ctx    = context.WithValue(context.Background(), ctxLogFieldsKey, fields)
	)

	AddFieldInContext(ctx, "key1", "value1")
	AddFieldInContext(ctx, "key2", "value2")

	assert.Equal(t, fields["key1"], "value1")
}

func TestAddErrorInContext(t *testing.T) {
	var (
		err error
		ctx = context.WithValue(context.Background(), ctxLogErrorsKey, &err)
	)

	AddErrorInContext(ctx, errors.New("eww1"))
	AddErrorInContext(ctx, errors.New("eww2"))

	assert.Equal(t, "eww2: eww1", err.Error())
}
