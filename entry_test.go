package log_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bep/log"
	"github.com/bep/log/handlers"
	"github.com/bep/log/handlers/memory"
	qt "github.com/frankban/quicktest"
)

func TestEntry_WithFields(t *testing.T) {
	h := memory.New()
	a := log.NewLogger(log.LoggerConfig{Handler: h, Level: log.InfoLevel}).WithLevel(log.InfoLevel)

	b := a.WithFields(log.Fields{{"foo", "bar"}})

	c := a.WithFields(log.Fields{{"foo", "hello"}, {"bar", "world"}})
	d := c.WithFields(log.Fields{{"baz", "jazz"}})
	qt.Assert(t, b.Fields, qt.DeepEquals, log.Fields{{"foo", "bar"}})
	qt.Assert(t, c.Fields, qt.DeepEquals, log.Fields{{"foo", "hello"}, {"bar", "world"}})
	qt.Assert(t, d.Fields, qt.DeepEquals, log.Fields{{"foo", "hello"}, {"bar", "world"}, {"baz", "jazz"}})

	c.Log(log.String("upload"))
	e := h.Entries[0]

	qt.Assert(t, "upload", qt.Equals, e.Message)
	qt.Assert(t, log.Fields{{"foo", "hello"}, {"bar", "world"}}, qt.DeepEquals, e.Fields)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
	qt.Assert(t, time.Now().IsZero(), qt.IsFalse)
}

func TestEntry_WithField(t *testing.T) {
	h := memory.New()
	a := log.NewLogger(log.LoggerConfig{Handler: h, Level: log.InfoLevel}).WithLevel(log.InfoLevel)
	b := a.WithField("foo", "baz").WithField("foo", "bar")
	b.Log(log.String("upload"))
	qt.Assert(t, a.Fields, qt.IsNil)
	qt.Assert(t, h.Entries[0].Fields, qt.DeepEquals, log.Fields{{"foo", "bar"}})
}

func TestEntry_WithError(t *testing.T) {
	a := log.NewLogger(log.LoggerConfig{Handler: handlers.Discard, Level: log.InfoLevel}).WithLevel(log.InfoLevel)
	b := a.WithError(fmt.Errorf("boom"))
	qt.Assert(t, a.Fields, qt.IsNil)
	qt.Assert(t, b.Fields, qt.DeepEquals, log.Fields{{"error", "boom"}})
}

func TestEntry_WithError_fields(t *testing.T) {
	a := log.NewLogger(log.LoggerConfig{Handler: handlers.Discard, Level: log.InfoLevel}).WithLevel(log.InfoLevel)
	b := a.WithError(errFields("boom"))
	qt.Assert(t, a.Fields, qt.IsNil)
	qt.Assert(t,

		b.Fields, qt.DeepEquals, log.Fields{
			{"error", "boom"},
			{"reason", "timeout"},
		})
}

func TestEntry_WithError_nil(t *testing.T) {
	a := log.NewLogger(log.LoggerConfig{Handler: handlers.Discard, Level: log.InfoLevel}).WithLevel(log.InfoLevel)
	b := a.WithError(nil)
	qt.Assert(t, a.Fields, qt.IsNil)
	qt.Assert(t, b.Fields, qt.IsNil)
}

func TestEntry_WithDuration(t *testing.T) {
	a := log.NewLogger(log.LoggerConfig{Handler: handlers.Discard, Level: log.InfoLevel}).WithLevel(log.InfoLevel)
	b := a.WithDuration(time.Second * 2)
	qt.Assert(t, b.Fields, qt.DeepEquals, log.Fields{{"duration", int64(2000)}})
}

type errFields string

func (ef errFields) Error() string {
	return string(ef)
}

func (ef errFields) Fields() log.Fields {
	return log.Fields{{"reason", "timeout"}}
}
