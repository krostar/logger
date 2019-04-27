package logger

import "github.com/pkg/errors"

// Config defines all the configurable options for the logger.
type Config struct {
	Verbosity string `json:"verbosity"  yaml:"verbosity"`
	Formatter string `json:"formatter"  yaml:"formatter"`
	WithColor bool   `json:"with-color" yaml:"with-color"`
	Output    string `json:"output"     yaml:"output"`
}

// SetDefault set sane default for logger's config.
func (c *Config) SetDefault() {
	c.Verbosity = LevelInfo.String()
	c.Formatter = "console"
	c.WithColor = true
	c.Output = "stdout"
}

// Validate makes sure the configuration is valid.
func (c *Config) Validate() error {
	if _, err := ParseLevel(c.Verbosity); err != nil {
		return errors.Wrapf(err, "unable to parse level %q", c.Verbosity)
	}

	switch c.Formatter {
	case "json":
	case "console":
	default:
		return errors.Errorf("unknown formatter %q", c.Formatter)
	}

	return nil
}
