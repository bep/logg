package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// assert interface compliance.
var (
	_ Interface = (*EntryFields)(nil)
	_ Interface = (*Entry)(nil)
)

// Now returns the current time.
var Now = time.Now

// EntryFields represents a single log entry.
type EntryFields struct {
	Logger *Logger `json:"-"`
	Fields Fields  `json:"-"`
}

// Entry holds a Entry with a Message, Timestamp and a Level.
// This is what is actually logged.
type Entry struct {
	*EntryFields
	Level     Level     `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

// FieldsDistinct returns a list of fields with duplicate names removed,
// keeping the last.
func (e *EntryFields) FieldsDistinct() Fields {
	return e.distinctFieldsLastByName()
}

func (e Entry) MarshalJSON() ([]byte, error) {
	fields := make(map[string]any)
	for _, f := range e.FieldsDistinct() {
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
func (e *EntryFields) Debug(s fmt.Stringer) {
	e.Logger.log(DebugLevel, e, s)
}

// Info level message.
func (e *EntryFields) Info(s fmt.Stringer) {
	e.Logger.log(InfoLevel, e, s)
}

// Warn level message.
func (e *EntryFields) Warn(s fmt.Stringer) {
	e.Logger.log(WarnLevel, e, s)
}

// Error level message.
func (e *EntryFields) Error(s fmt.Stringer) {
	e.Logger.log(ErrorLevel, e, s)
}

// Fatal level message, followed by an exit.
func (e *EntryFields) Fatal(s fmt.Stringer) {
	e.Logger.log(FatalLevel, e, s)
	os.Exit(1)
}

// distinctFieldsLastByName returns the fields with duplicate names removed,
// keeping the rightmost field (last) with a given name.
func (e *EntryFields) distinctFieldsLastByName() Fields {
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
			// Insert first.
			fields = append(fields, Field{})
			copy(fields[1:], fields[:])
			fields[0] = f
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
		EntryFields: e,
		Level:       level,
		Message:     msg,
		Timestamp:   Now(),
	}
}
