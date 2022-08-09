package logg_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bep/logg"
	"github.com/bep/logg/handlers"
	"github.com/bep/logg/handlers/memory"
	qt "github.com/frankban/quicktest"
)

func TestEntry_WithFields(t *testing.T) {
	h := memory.New()
	a := logg.NewLogger(logg.LoggerConfig{Handler: h, Level: logg.LevelInfo}).WithLevel(logg.LevelInfo)

	b := a.WithFields(logg.Fields{{"foo", "bar"}})

	c := a.WithFields(logg.Fields{{"foo", "hello"}, {"bar", "world"}})
	d := c.WithFields(logg.Fields{{"baz", "jazz"}})
	qt.Assert(t, b.Fields, qt.DeepEquals, logg.Fields{{"foo", "bar"}})
	qt.Assert(t, c.Fields, qt.DeepEquals, logg.Fields{{"foo", "hello"}, {"bar", "world"}})
	qt.Assert(t, d.Fields, qt.DeepEquals, logg.Fields{{"foo", "hello"}, {"bar", "world"}, {"baz", "jazz"}})

	c.Log(logg.String("upload"))
	e := h.Entries[0]

	qt.Assert(t, "upload", qt.Equals, e.Message)
	qt.Assert(t, logg.Fields{{"foo", "hello"}, {"bar", "world"}}, qt.DeepEquals, e.Fields)
	qt.Assert(t, logg.LevelInfo, qt.Equals, e.Level)
	qt.Assert(t, time.Now().IsZero(), qt.IsFalse)
}

func TestEntry_WithManyFieldsWithSameName(t *testing.T) {
	h := memory.New()
	a := logg.NewLogger(logg.LoggerConfig{Handler: h, Level: logg.LevelInfo}).WithLevel(logg.LevelInfo)

	b := a.WithFields(logg.Fields{{"foo", "bar"}})

	for i := 0; i < 100; i++ {
		b = b.WithFields(logg.Fields{{"foo", "bar"}})
	}

	b.Log(logg.String("upload"))
	e := h.Entries[0]

	qt.Assert(t, "upload", qt.Equals, e.Message)
	qt.Assert(t, logg.Fields{{"foo", "bar"}}, qt.DeepEquals, e.Fields)

}

func TestEntry_WithField(t *testing.T) {
	h := memory.New()
	a := logg.NewLogger(logg.LoggerConfig{Handler: h, Level: logg.LevelInfo}).WithLevel(logg.LevelInfo)
	b := a.WithField("foo", "baz").WithField("foo", "bar")
	b.Log(logg.String("upload"))
	qt.Assert(t, a.Fields, qt.IsNil)
	qt.Assert(t, h.Entries[0].Fields, qt.DeepEquals, logg.Fields{{"foo", "bar"}})
}

func TestEntry_WithError(t *testing.T) {
	a := logg.NewLogger(logg.LoggerConfig{Handler: handlers.Discard, Level: logg.LevelInfo}).WithLevel(logg.LevelInfo)
	b := a.WithError(fmt.Errorf("boom"))
	qt.Assert(t, a.Fields, qt.IsNil)
	qt.Assert(t, b.Fields, qt.DeepEquals, logg.Fields{{"error", "boom"}})
}

func TestEntry_WithError_fields(t *testing.T) {
	a := logg.NewLogger(logg.LoggerConfig{Handler: handlers.Discard, Level: logg.LevelInfo}).WithLevel(logg.LevelInfo)
	b := a.WithError(errFields("boom"))
	qt.Assert(t, a.Fields, qt.IsNil)
	qt.Assert(t,

		b.Fields, qt.DeepEquals, logg.Fields{
			{"error", "boom"},
			{"reason", "timeout"},
		})
}

func TestEntry_WithError_nil(t *testing.T) {
	a := logg.NewLogger(logg.LoggerConfig{Handler: handlers.Discard, Level: logg.LevelInfo}).WithLevel(logg.LevelInfo)
	b := a.WithError(nil)
	qt.Assert(t, a.Fields, qt.IsNil)
	qt.Assert(t, b.Fields, qt.IsNil)
}

func TestEntry_WithDuration(t *testing.T) {
	a := logg.NewLogger(logg.LoggerConfig{Handler: handlers.Discard, Level: logg.LevelInfo}).WithLevel(logg.LevelInfo)
	b := a.WithDuration(time.Second * 2)
	qt.Assert(t, b.Fields, qt.DeepEquals, logg.Fields{{"duration", int64(2000)}})
}

type errFields string

func (ef errFields) Error() string {
	return string(ef)
}

func (ef errFields) Fields() logg.Fields {
	return logg.Fields{{"reason", "timeout"}}
}
