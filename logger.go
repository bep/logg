package log

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/bep/clocks"
)

// assert interface compliance.
var _ Leveler = (*logger)(nil)

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

// NewStringFunc returns a StringFunc which implements Stringer.
func NewStringFunc(f func() string) StringFunc {
	return StringFunc(f)
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

// NoopHandler is a no-op handler that discards all log messages.
var NoopHandler = HandlerFunc(func(e *Entry) error {
	return nil
})

// LoggerConfig is the configuration used to create a logger.
type LoggerConfig struct {
	// Level is the minimum level to log at.
	// If not set, defaults to InfoLevel.
	Level Level

	// Handler is the log handler to use.
	Handler Handler

	// Clock is the clock to use for timestamps.
	// If not set, the system clock is used.
	Clock Clock
}

// New returns a new logger.
func NewLogger(cfg LoggerConfig) Leveler {
	if cfg.Handler == nil {
		panic("handler cannot be nil")
	}

	if cfg.Level <= 0 || cfg.Level > FatalLevel {
		panic("log level is out of range")
	}

	if cfg.Clock == nil {
		cfg.Clock = clocks.System()
	}

	if cfg.Level == 0 {
		cfg.Level = InfoLevel
	}

	return &logger{
		Handler: cfg.Handler,
		Level:   cfg.Level,
		Clock:   cfg.Clock,
	}
}

// logger represents a logger with configurable Level and Handler.
type logger struct {
	Handler Handler
	Level   Level
	Clock   Clock
}

// Clock provides the current time.
type Clock interface {
	Now() time.Time
}

// WithLevel returns a new entry with `level` set.
func (l *logger) WithLevel(level Level) *EntryFields {
	return NewEntry(l).WithLevel(level)
}

// WithFields returns a new entry with `fields` set.
func (l *logger) WithFields(fields Fielder) *EntryFields {
	return NewEntry(l).WithFields(fields.Fields())
}

// WithField returns a new entry with the `key` and `value` set.
//
// Note that the `key` should not have spaces in it - use camel
// case or underscores
func (l *logger) WithField(key string, value any) *EntryFields {
	return NewEntry(l).WithField(key, value)
}

// WithDuration returns a new entry with the "duration" field set
// to the given duration in milliseconds.
func (l *logger) WithDuration(d time.Duration) *EntryFields {
	return NewEntry(l).WithDuration(d)
}

// WithError returns a new entry with the "error" set to `err`.
func (l *logger) WithError(err error) *EntryFields {
	return NewEntry(l).WithError(err)
}

// log the message, invoking the handler. We clone the entry here
// to bypass the overhead in Entry methods when the level is not
// met.
func (l *logger) log(e *EntryFields, s fmt.Stringer) {
	if e.Level < l.Level {
		return
	}

	if err := l.Handler.HandleLog(e.finalize(s.String())); err != nil {
		stdlog.Printf("error logging: %s", err)
	}

	if e.Level == FatalLevel {
		os.Exit(1)
	}
}
