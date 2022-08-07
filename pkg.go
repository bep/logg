package log

import (
	"fmt"
	"time"
)

// singletons ftw?
var Log Interface = &Logger{
	Handler: HandlerFunc(handleStdLog),
	Level:   InfoLevel,
}

// SetHandler sets the handler. This is not thread-safe.
// The default handler outputs to the stdlib log.
func SetHandler(h Handler) {
	if logger, ok := Log.(*Logger); ok {
		logger.Handler = h
	}
}

// SetLevel sets the log level. This is not thread-safe.
func SetLevel(l Level) {
	if logger, ok := Log.(*Logger); ok {
		logger.Level = l
	}
}

// SetLevelFromString sets the log level from a string, panicing when invalid. This is not thread-safe.
func SetLevelFromString(s string) {
	if logger, ok := Log.(*Logger); ok {
		logger.Level = MustParseLevel(s)
	}
}

// WithFields returns a new entry with `fields` set.
func WithFields(fields Fielder) *EntryFields {
	return Log.WithFields(fields)
}

// WithField returns a new entry with the `key` and `value` set.
func WithField(key string, value any) *EntryFields {
	return Log.WithField(key, value)
}

// WithDuration returns a new entry with the "duration" field set
// to the given duration in milliseconds.
func WithDuration(d time.Duration) *EntryFields {
	return Log.WithDuration(d)
}

// WithError returns a new entry with the "error" set to `err`.
func WithError(err error) *EntryFields {
	return Log.WithError(err)
}

// Debug level message.
func Debug(s fmt.Stringer) {
	Log.Debug(s)
}

// Info level message.
func Info(s fmt.Stringer) {
	Log.Info(s)
}

// Warn level message.
func Warn(s fmt.Stringer) {
	Log.Warn(s)
}

// Error level message.
func Error(s fmt.Stringer) {
	Log.Error(s)
}

// Fatal level message, followed by an exit.
func Fatal(s fmt.Stringer) {
	Log.Fatal(s)
}
