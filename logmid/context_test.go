package logmid

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AddFieldInContext(t *testing.T) {
	var (
		fields = make(map[string]interface{})
		ctx    = context.WithValue(context.Background(), ctxLogFieldsKey, fields)
	)

	AddFieldInContext(ctx, "key1", "value1")
	AddFieldInContext(ctx, "key2", "value2")

	assert.Equal(t, fields["key1"], "value1")
}

func Test_AddErrorInContext(t *testing.T) {
	var (
		err error
		ctx = context.WithValue(context.Background(), ctxLogErrorsKey, &err)
	)

	AddErrorInContext(ctx, errors.New("eww1"))
	AddErrorInContext(ctx, errors.New("eww2"))

	require.Error(t, err)
	assert.Equal(t, "eww2: eww1", err.Error())
}
