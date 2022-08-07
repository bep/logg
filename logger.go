package log

import (
	"fmt"

	"time"
)

// assert interface compliance.
var _ Interface = (*Logger)(nil)

// String implements fmt.Stringer and can be used directly in
// the log methods.
type String string

// StringFunc is a function that returns a string.
// It also implements the fmt.Stringer interface and
// can therefore be used as argument to the log methods.
type StringFunc func() string

func (f StringFunc) String() string {
	return f()
}

func (s String) String() string {
	return string(s)
}

// Fielder is an interface for providing fields to custom types.
type Fielder interface {
	Fields() Fields
}

// Field holds a named value.
type Field struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

// Fields represents a slice of entry level data used for structured logging.
type Fields []Field

// Fields implements Fielder.
func (f Fields) Fields() Fields {
	return f
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as
// log handlers. If f is a function with the appropriate signature,
// HandlerFunc(f) is a Handler object that calls f.
type HandlerFunc func(*Entry) error

// HandleLog calls f(e).
func (f HandlerFunc) HandleLog(e *Entry) error {
	return f(e)
}

// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services. See the "handlers"
// directory for implementations.
//
// It is left up to Handlers to implement thread-safety.
type Handler interface {
	HandleLog(*Entry) error
}

// Logger represents a logger with configurable Level and Handler.
type Logger struct {
	Handler Handler
	Level   Level
}

// WithFields returns a new entry with `fields` set.
func (l *Logger) WithFields(fields Fielder) *EntryFields {
	return NewEntry(l).WithFields(fields.Fields())
}

// WithField returns a new entry with the `key` and `value` set.
//
// Note that the `key` should not have spaces in it - use camel
// case or underscores
func (l *Logger) WithField(key string, value any) *EntryFields {
	return NewEntry(l).WithField(key, value)
}

// WithDuration returns a new entry with the "duration" field set
// to the given duration in milliseconds.
func (l *Logger) WithDuration(d time.Duration) *EntryFields {
	return NewEntry(l).WithDuration(d)
}

// WithError returns a new entry with the "error" set to `err`.
func (l *Logger) WithError(err error) *EntryFields {
	return NewEntry(l).WithError(err)
}

// Debug level message.
func (l *Logger) Debug(s fmt.Stringer) {
	if DebugLevel < l.Level {
		return
	}
	NewEntry(l).Debug(s)
}

// Info level message.
func (l *Logger) Info(s fmt.Stringer) {
	if InfoLevel < l.Level {
		return
	}
	NewEntry(l).Info(s)
}

// Warn level message.
func (l *Logger) Warn(s fmt.Stringer) {
	if WarnLevel < l.Level {
		return
	}
	NewEntry(l).Warn(s)
}

// Error level message.
func (l *Logger) Error(s fmt.Stringer) {
	if ErrorLevel < l.Level {
		return
	}
	NewEntry(l).Error(s)
}

// Fatal level message, followed by an exit.
func (l *Logger) Fatal(s fmt.Stringer) {
	if FatalLevel < l.Level {
		return
	}
	NewEntry(l).Fatal(s)
}

// log the message, invoking the handler. We clone the entry here
// to bypass the overhead in Entry methods when the level is not
// met.
func (l *Logger) log(level Level, e *EntryFields, s fmt.Stringer) {
	if level < l.Level {
		return
	}

	if err := l.Handler.HandleLog(e.finalize(level, s.String())); err != nil {
		panic(fmt.Sprintf("log: error invoking handler: %v", err))
	}
}
