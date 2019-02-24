package logger

// Config defines all the configurable options for the logger.
type Config struct {
	Verbosity string `json:"verbosity"  yaml:"verbosity"`
	Formatter string `json:"formatter"  yaml:"formatter"`
	WithColor bool   `json:"with-color" yaml:"with-color"`
	Output    string `json:"output"     yaml:"output"`
	OutputErr string `json:"output-err" yaml:"output-err"`
}

// SetDefault set sane default for logger's config.
func (c *Config) SetDefault() {
	c.Verbosity = LevelInfo.String()
	c.Formatter = "console"
	c.WithColor = true
	c.Output = "stdout"
	c.OutputErr = "stderr"
}
