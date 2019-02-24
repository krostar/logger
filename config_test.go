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
			OutputErr: "stderr",
		}
	)

	cfg.SetDefault()

	assert.Equal(t, expectedCfg, cfg)
}
