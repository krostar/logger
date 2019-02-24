package logger

// InMemory defines a memory logger.
// It is designed for tests purposes only.
type InMemory struct {
	Level   Level
	Fields  map[string]interface{}
	Entries []InMemoryEntry
}

// InMemoryEntry stores a log entry for the Memory logger.
type InMemoryEntry struct {
	Level  Level
	Format string
	Args   []interface{}
}

// NewInMemory returns a logger that stores log in memory.
func NewInMemory(lvl Level) *InMemory {
	return &InMemory{
		Level:   lvl,
		Fields:  make(map[string]interface{}),
		Entries: []InMemoryEntry{},
	}
}

// Reset clears all fields and entries.
func (n *InMemory) Reset(lvl Level) {
	for key := range n.Fields {
		delete(n.Fields, key)
	}

	n.Entries = nil
	n.Level = lvl
}

// RedirectStdLog is a no-op function, for now.
func (n *InMemory) RedirectStdLog(at Level) (func(), error) {
	return nil, nil
}

// SetLevel implements Logger for Memory.
func (n *InMemory) SetLevel(lvl Level) error {
	n.Level = lvl
	return nil
}

// Debug implements Logger for Memory.
func (n *InMemory) Debug(args ...interface{}) {
	if n.Level <= LevelDebug {
		n.Entries = append(n.Entries, InMemoryEntry{LevelDebug, "", args})
	}
}

// Debugf implements Logger for Memory.
func (n *InMemory) Debugf(format string, args ...interface{}) {
	if n.Level <= LevelDebug {
		n.Entries = append(n.Entries, InMemoryEntry{LevelDebug, format, args})
	}
}

// Info implements Logger for Memory.
func (n *InMemory) Info(args ...interface{}) {
	if n.Level <= LevelInfo {
		n.Entries = append(n.Entries, InMemoryEntry{LevelInfo, "", args})
	}
}

// Infof implements Logger for Memory.
func (n *InMemory) Infof(format string, args ...interface{}) {
	if n.Level <= LevelInfo {
		n.Entries = append(n.Entries, InMemoryEntry{LevelInfo, format, args})
	}
}

// Warn implements Logger for Memory.
func (n *InMemory) Warn(args ...interface{}) {
	if n.Level <= LevelWarn {
		n.Entries = append(n.Entries, InMemoryEntry{LevelWarn, "", args})
	}
}

// Warnf implements Logger for Memory.
func (n *InMemory) Warnf(format string, args ...interface{}) {
	if n.Level <= LevelWarn {
		n.Entries = append(n.Entries, InMemoryEntry{LevelWarn, format, args})
	}
}

// Error implements Logger for Memory.
func (n *InMemory) Error(args ...interface{}) {
	if n.Level <= LevelError {
		n.Entries = append(n.Entries, InMemoryEntry{LevelError, "", args})
	}
}

// Errorf implements Logger for Memory.
func (n *InMemory) Errorf(format string, args ...interface{}) {
	if n.Level <= LevelError {
		n.Entries = append(n.Entries, InMemoryEntry{LevelError, format, args})
	}
}

// WithField implements Logger for Memory.
func (n *InMemory) WithField(key string, value interface{}) Logger {
	n.Fields[key] = value
	return n
}

// WithFields implements Logger for Memory.
func (n *InMemory) WithFields(fields map[string]interface{}) Logger {
	for key, value := range fields {
		n.Fields[key] = value
	}
	return n
}

// WithError implements Logger for Memory.
func (n *InMemory) WithError(err error) Logger {
	n.Fields[FieldErrorKey] = err
	return n
}
