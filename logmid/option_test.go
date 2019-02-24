package logmid

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/krostar/logger"
)

func TestWithCallback(t *testing.T) {
	var opts Options

	WithCallback(func(r *http.Request) {})(&opts)

	assert.Len(t, opts.onRequest, 1)
}

func TestWithLogLevelFunc(t *testing.T) {
	var opts Options

	WithLogLevelFunc(func(r *http.Request) logger.Level { return logger.LevelDebug })(&opts)

	assert.NotNil(t, opts.logAtLevel)
}
