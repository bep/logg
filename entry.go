package logg

import (
	"fmt"
	"strings"
	"time"
)

// assert interface compliance.
var (
	_ LevelLogger = (*Entry)(nil)
)

// Entry represents a single log entry at a given log level.
type Entry struct {
	logger *logger

	Level     Level     `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Fields    Fields    `json:"fields,omitempty"`
	Message   string    `json:"message"`
}

// NewEntry returns a new entry for `log`.
func NewEntry(log *logger) *Entry {
	return &Entry{
		logger: log,
	}
}

func (e Entry) WithLevel(level Level) *Entry {
	e.Level = level
	return &e
}

func (e *Entry) WithFields(fielder Fielder) *Entry {
	if e.isLevelDisabled() {
		return e
	}
	x := *e
	x.Fields = append(x.Fields, fielder.Fields()...)
	return &x
}

func (e *Entry) WithField(key string, value any) *Entry {
	if e.isLevelDisabled() {
		return e
	}
	return e.WithFields(Fields{{key, value}})
}

func (e *Entry) WithDuration(d time.Duration) *Entry {
	if e.isLevelDisabled() {
		return e
	}
	return e.WithField("duration", d.Milliseconds())
}

// WithError returns a new entry with the "error" set to `err`.
//
// The given error may implement .Fielder, if it does the method
// will add all its `.Fields()` into the returned entry.
func (e *Entry) WithError(err error) *Entry {
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

func (e *Entry) isLevelDisabled() bool {
	return e.Level < e.logger.Level
}

// Log a message at the given level.
func (e *Entry) Log(s fmt.Stringer) {
	e.logger.log(e, s)
}

// Clone returns a new Entry with the same fields.
func (e *Entry) Clone() *Entry {
	x := *e
	x.Fields = make(Fields, len(e.Fields))
	copy(x.Fields, e.Fields)
	return &x
}

func (e *Entry) reset() {
	e.logger = nil
	e.Level = 0
	e.Fields = e.Fields[:0]
	e.Message = ""
	e.Timestamp = time.Time{}
}

// finalize populates dst with Level and  Fields merged from e and Message and Timestamp set.
func (e *Entry) finalize(dst *Entry, msg string) {
	dst.Message = msg
	dst.Timestamp = e.logger.Clock.Now()
	dst.Level = e.Level

	// There mau be fields logged with the same name, keep the latest.
	for i := len(e.Fields) - 1; i >= 0; i-- {
		f := e.Fields[i]
		var seen bool
		for _, f2 := range dst.Fields {
			if f.Name == f2.Name {
				seen = true
				break
			}
		}
		if !seen {
			// Insert first.
			dst.Fields = append(dst.Fields, Field{})
			copy(dst.Fields[1:], dst.Fields[:])
			dst.Fields[0] = f
		}
	}
}
