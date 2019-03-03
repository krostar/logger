package logger

// Noop defines a no-operation logger.
type Noop struct{}

// SetLevel implements Logger for Noop.
func (n *Noop) SetLevel(Level) error { return nil }

// Debug implements Logger for Noop.
func (*Noop) Debug(...interface{}) {}

// Debugf implements Logger for Noop.
func (*Noop) Debugf(string, ...interface{}) {}

// Info implements Logger for Noop.
func (*Noop) Info(...interface{}) {}

// Infof implements Logger for Noop.
func (*Noop) Infof(string, ...interface{}) {}

// Warn implements Logger for Noop.
func (*Noop) Warn(...interface{}) {}

// Warnf implements Logger for Noop.
func (*Noop) Warnf(string, ...interface{}) {}

// Error implements Logger for Noop.
func (*Noop) Error(...interface{}) {}

// Errorf implements Logger for Noop.
func (*Noop) Errorf(string, ...interface{}) {}

// WithField implements Logger for Noop.
func (n *Noop) WithField(string, interface{}) Logger { return &Noop{} }

// WithFields implements Logger for Noop.
func (n *Noop) WithFields(map[string]interface{}) Logger { return &Noop{} }

// WithError implements Logger for Noop.
func (n *Noop) WithError(error) Logger { return &Noop{} }
