package log

import (
	stdlog "log"
	"time"
)

// assert interface compliance.
var _ Interface = (*Logger)(nil)

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
func (l *Logger) Debug(msg string) {
	NewEntry(l).Debug(msg)
}

// Info level message.
func (l *Logger) Info(msg string) {
	NewEntry(l).Info(msg)
}

// Warn level message.
func (l *Logger) Warn(msg string) {
	NewEntry(l).Warn(msg)
}

// Error level message.
func (l *Logger) Error(msg string) {
	NewEntry(l).Error(msg)
}

// Fatal level message, followed by an exit.
func (l *Logger) Fatal(msg string) {
	NewEntry(l).Fatal(msg)
}

// Debugf level formatted message.
func (l *Logger) Debugf(msg string, v ...any) {
	NewEntry(l).Debugf(msg, v...)
}

// Infof level formatted message.
func (l *Logger) Infof(msg string, v ...any) {
	NewEntry(l).Infof(msg, v...)
}

// Warnf level formatted message.
func (l *Logger) Warnf(msg string, v ...any) {
	NewEntry(l).Warnf(msg, v...)
}

// Errorf level formatted message.
func (l *Logger) Errorf(msg string, v ...any) {
	NewEntry(l).Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func (l *Logger) Fatalf(msg string, v ...any) {
	NewEntry(l).Fatalf(msg, v...)
}

// log the message, invoking the handler. We clone the entry here
// to bypass the overhead in Entry methods when the level is not
// met.
func (l *Logger) log(level Level, e *EntryFields, msg string) {
	if level < l.Level {
		return
	}

	if err := l.Handler.HandleLog(e.finalize(level, msg)); err != nil {
		stdlog.Printf("error logging: %s", err)
	}
}
