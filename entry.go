package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// assert interface compliance.
var _ Interface = (*EntryFields)(nil)

// Now returns the current time.
var Now = time.Now

// EntryFields represents a single log entry.
type EntryFields struct {
	Logger *Logger `json:"-"`
	Fields Fields  `json:"-"`
	start  time.Time
}

// Entry holds a Entry with a Message, Timestamp and a Level.
// This is what is actually logged.
type Entry struct {
	*EntryFields
	FieldsUnique Fields    `json:"-"`
	Level        Level     `json:"level"`
	Timestamp    time.Time `json:"timestamp"`
	Message      string    `json:"message"`
}

func (e Entry) MarshalJSON() ([]byte, error) {
	fields := make(map[string]any)
	for _, f := range e.FieldsUnique {
		fields[f.Name] = f.Value
	}

	type EntryAlias Entry
	return json.Marshal(&struct {
		Fields map[string]any `json:"fields"`
		EntryAlias
	}{
		Fields:     fields,
		EntryAlias: (EntryAlias)(e),
	})
}

// NewEntry returns a new entry for `log`.
func NewEntry(log *Logger) *EntryFields {
	return &EntryFields{
		Logger: log,
	}
}

// WithFields returns a new entry with `fields` set.
func (e EntryFields) WithFields(fielder Fielder) *EntryFields {
	e.Fields = append(e.Fields, fielder.Fields()...)
	return &e
}

// WithField returns a new entry with the `key` and `value` set.
func (e *EntryFields) WithField(key string, value any) *EntryFields {
	return e.WithFields(Fields{{key, value}})
}

// WithDuration returns a new entry with the "duration" field set
// to the given duration in milliseconds.
func (e *EntryFields) WithDuration(d time.Duration) *EntryFields {
	return e.WithField("duration", d.Milliseconds())
}

// WithError returns a new entry with the "error" set to `err`.
//
// The given error may implement .Fielder, if it does the method
// will add all its `.Fields()` into the returned entry.
func (e *EntryFields) WithError(err error) *EntryFields {
	if err == nil {
		return e
	}

	ctx := e.WithField("error", err.Error())

	if s, ok := err.(stackTracer); ok {
		frame := s.StackTrace()[0]

		name := fmt.Sprintf("%n", frame)
		file := fmt.Sprintf("%+s", frame)
		line := fmt.Sprintf("%d", frame)

		parts := strings.Split(file, "\n\t")
		if len(parts) > 1 {
			file = parts[1]
		}

		ctx = ctx.WithField("source", fmt.Sprintf("%s: %s:%s", name, file, line))
	}

	if f, ok := err.(Fielder); ok {
		ctx = ctx.WithFields(f)
	}

	return ctx
}

// Debug level message.
func (e *EntryFields) Debug(msg string) {
	e.Logger.log(DebugLevel, e, msg)
}

// Info level message.
func (e *EntryFields) Info(msg string) {
	e.Logger.log(InfoLevel, e, msg)
}

// Warn level message.
func (e *EntryFields) Warn(msg string) {
	e.Logger.log(WarnLevel, e, msg)
}

// Error level message.
func (e *EntryFields) Error(msg string) {
	e.Logger.log(ErrorLevel, e, msg)
}

// Fatal level message, followed by an exit.
func (e *EntryFields) Fatal(msg string) {
	e.Logger.log(FatalLevel, e, msg)
	os.Exit(1)
}

// Debugf level formatted message.
func (e *EntryFields) Debugf(msg string, v ...any) {
	e.Debug(fmt.Sprintf(msg, v...))
}

// Infof level formatted message.
func (e *EntryFields) Infof(msg string, v ...any) {
	e.Info(fmt.Sprintf(msg, v...))
}

// Warnf level formatted message.
func (e *EntryFields) Warnf(msg string, v ...any) {
	e.Warn(fmt.Sprintf(msg, v...))
}

// Errorf level formatted message.
func (e *EntryFields) Errorf(msg string, v ...any) {
	e.Error(fmt.Sprintf(msg, v...))
}

// Fatalf level formatted message, followed by an exit.
func (e *EntryFields) Fatalf(msg string, v ...any) {
	e.Fatal(fmt.Sprintf(msg, v...))
}

// mergedFields returns the fields list collapsed into a single slice.
func (e *EntryFields) mergedFields() Fields {
	fields := make(Fields, 0, len(e.Fields))
	for i := len(e.Fields) - 1; i >= 0; i-- {
		f := e.Fields[i]
		var seen bool
		for _, f2 := range fields {
			if f.Name == f2.Name {
				seen = true
				break
			}
		}
		if !seen {
			fields = append(fields, f)
		}
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}

// finalize returns a copy of the Entry with Fields merged.
func (e *EntryFields) finalize(level Level, msg string) *Entry {
	return &Entry{
		EntryFields:  e,
		FieldsUnique: e.mergedFields(),
		Level:        level,
		Message:      msg,
		Timestamp:    Now(),
	}
}
