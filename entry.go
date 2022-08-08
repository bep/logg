package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// assert interface compliance.
var (
	_ LevelLogger = (*EntryFields)(nil)
	_ LevelLogger = (*Entry)(nil)
)

// EntryFields represents a single log entry at a given log level.
type EntryFields struct {
	Logger *logger `json:"-"`
	Fields Fields  `json:"-"`
	Level  Level   `json:"level"`
}

// Entry holds a Entry with a Message and a Timestamp.
// This is what is actually logged.
type Entry struct {
	*EntryFields
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
func NewEntry(log *logger) *EntryFields {
	return &EntryFields{
		Logger: log,
	}
}

func (e EntryFields) WithLevel(level Level) *EntryFields {
	e.Level = level
	return &e
}

func (e *EntryFields) WithFields(fielder Fielder) *EntryFields {
	if e.isLevelDisabled() {
		return e
	}
	x := *e
	x.Fields = append(x.Fields, fielder.Fields()...)
	return &x
}

func (e *EntryFields) WithField(key string, value any) *EntryFields {
	if e.isLevelDisabled() {
		return e
	}
	return e.WithFields(Fields{{key, value}})
}

func (e *EntryFields) WithDuration(d time.Duration) *EntryFields {
	if e.isLevelDisabled() {
		return e
	}
	return e.WithField("duration", d.Milliseconds())
}

// WithError returns a new entry with the "error" set to `err`.
//
// The given error may implement .Fielder, if it does the method
// will add all its `.Fields()` into the returned entry.
func (e *EntryFields) WithError(err error) *EntryFields {
	if err == nil || e.isLevelDisabled() {
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

func (e *EntryFields) isLevelDisabled() bool {
	return e.Level < e.Logger.Level
}

// Log a message at the given level.
func (e *EntryFields) Log(s fmt.Stringer) {
	e.Logger.log(e, s)
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
func (e *EntryFields) finalize(msg string) *Entry {
	return &Entry{
		EntryFields: e,
		Message:     msg,
		Timestamp:   e.Logger.Clock.Now(),
	}
}
