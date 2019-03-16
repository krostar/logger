package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_SetDefault(t *testing.T) {
	var (
		cfg         Config
		expectedCfg = Config{
			Verbosity: LevelInfo.String(),
			Formatter: "console",
			WithColor: true,
			Output:    "stdout",
		}
	)

	cfg.SetDefault()

	assert.Equal(t, expectedCfg, cfg)
}

func TestConfig_Validate(t *testing.T) {
	t.Run("default is valid", func(t *testing.T) {
		var cfg Config
		cfg.SetDefault()
		assert.NoError(t, cfg.Validate())
	})

	t.Run("verbosity fail", func(t *testing.T) {
		var cfg Config
		cfg.SetDefault()

		cfg.Verbosity = "boum"
		assert.Error(t, cfg.Validate())
	})

	t.Run("formatter fail", func(t *testing.T) {
		var cfg Config
		cfg.SetDefault()

		cfg.Formatter = "boum"
		assert.Error(t, cfg.Validate())
	})
}
