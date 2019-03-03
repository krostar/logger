package logger

// InMemory defines a memory logger.
// It is designed for tests purposes only.
type InMemory struct {
	parent  *InMemory
	fields  map[string]interface{}
	level   Level
	Entries []InMemoryEntry
}

// InMemoryEntry stores a log entry for the Memory logger.
type InMemoryEntry struct {
	Level  Level
	Format string
	Args   []interface{}
	Fields map[string]interface{}
}

// NewInMemory returns a logger that stores log in memory.
func NewInMemory(lvl Level) *InMemory {
	return &InMemory{
		level:   lvl,
		fields:  make(map[string]interface{}),
		Entries: []InMemoryEntry{},
	}
}

// Reset clears all fields and entries.
func (n *InMemory) Reset() {
	for key := range n.fields {
		delete(n.fields, key)
	}

	n.Entries = []InMemoryEntry{}
}

// SetLevel implements Logger for Memory.
func (n *InMemory) SetLevel(lvl Level) error {
	n.level = lvl
	return nil
}

// Debug implements Logger for Memory.
func (n *InMemory) Debug(args ...interface{}) { n.log(nil, LevelDebug, "", args) }

// Debugf implements Logger for Memory.
func (n *InMemory) Debugf(format string, args ...interface{}) { n.log(nil, LevelDebug, format, args) }

// Info implements Logger for Memory.
func (n *InMemory) Info(args ...interface{}) { n.log(nil, LevelInfo, "", args) }

// Infof implements Logger for Memory.
func (n *InMemory) Infof(format string, args ...interface{}) { n.log(nil, LevelInfo, format, args) }

// Warn implements Logger for Memory.
func (n *InMemory) Warn(args ...interface{}) { n.log(nil, LevelWarn, "", args) }

// Warnf implements Logger for Memory.
func (n *InMemory) Warnf(format string, args ...interface{}) { n.log(nil, LevelWarn, format, args) }

// Error implements Logger for Memory.
func (n *InMemory) Error(args ...interface{}) { n.log(nil, LevelError, "", args) }

// Errorf implements Logger for Memory.
func (n *InMemory) Errorf(format string, args ...interface{}) { n.log(nil, LevelError, format, args) }

// WithField implements Logger for Memory.
func (n *InMemory) WithField(key string, value interface{}) Logger {
	child := NewInMemory(n.level)
	child.parent = n

	child.fields[key] = value
	return child
}

// WithFields implements Logger for Memory.
func (n *InMemory) WithFields(fields map[string]interface{}) Logger {
	child := NewInMemory(n.level)
	child.parent = n

	for key, value := range fields {
		child.fields[key] = value
	}
	return child
}

// WithError implements Logger for Memory.
func (n *InMemory) WithError(err error) Logger {
	return n.WithField(FieldErrorKey, err)
}

func (n *InMemory) log(childFields map[string]interface{}, lvl Level, format string, args []interface{}) {
	var fields = make(map[string]interface{})

	for key, value := range childFields {
		fields[key] = value
	}
	for key, value := range n.fields {
		fields[key] = value
	}

	if lvl >= n.level {
		n.Entries = append(n.Entries, InMemoryEntry{
			Level:  lvl,
			Format: format,
			Args:   args,
			Fields: fields,
		})
	}

	if n.parent != nil {
		n.parent.log(fields, lvl, format, args)
	}
}
