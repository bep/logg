package log

import (
	"fmt"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestEntry_WithFields(t *testing.T) {
	a := NewEntry(nil)
	qt.Assert(t, a.Fields, qt.IsNil)

	b := a.WithFields(Fields{{"foo", "bar"}})
	qt.Assert(t, a.mergedFields(), qt.IsNil)

	c := a.WithFields(Fields{{"foo", "hello"}, {"bar", "world"}})
	d := c.WithFields(Fields{{"baz", "jazz"}})
	qt.Assert(t, b.mergedFields(), qt.DeepEquals, Fields{{"foo", "bar"}})
	qt.Assert(t, c.mergedFields(), qt.DeepEquals, Fields{{"foo", "hello"}, {"bar", "world"}})
	qt.Assert(t, d.mergedFields(), qt.DeepEquals, Fields{{"foo", "hello"}, {"bar", "world"}, {"baz", "jazz"}})

	e := c.finalize(InfoLevel, "upload")
	qt.Assert(t, "upload", qt.Equals, e.Message)
	qt.Assert(t, Fields{{"foo", "hello"}, {"bar", "world"}}, qt.DeepEquals, e.Fields)
	qt.Assert(t, InfoLevel, qt.Equals, e.Level)
	qt.Assert(t, time.Now().IsZero(), qt.IsFalse)
}

func TestEntry_WithField(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithField("foo", "baz").WithField("foo", "bar")
	qt.Assert(t, a.mergedFields(), qt.IsNil)
	qt.Assert(t, b.mergedFields(), qt.DeepEquals, Fields{{"foo", "bar"}})
}

func TestEntry_WithError(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithError(fmt.Errorf("boom"))
	qt.Assert(t, a.mergedFields(), qt.IsNil)
	qt.Assert(t, b.mergedFields(), qt.DeepEquals, Fields{{"error", "boom"}})
}

func TestEntry_WithError_fields(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithError(errFields("boom"))
	qt.Assert(t, a.mergedFields(), qt.IsNil)
	qt.Assert(t,

		b.mergedFields(), qt.DeepEquals, Fields{
			{"error", "boom"},
			{"reason", "timeout"},
		})

}

func TestEntry_WithError_nil(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithError(nil)
	qt.Assert(t, a.mergedFields(), qt.IsNil)
	qt.Assert(t, b.mergedFields(), qt.IsNil)
}

func TestEntry_WithDuration(t *testing.T) {
	a := NewEntry(nil)
	b := a.WithDuration(time.Second * 2)
	qt.Assert(t, b.mergedFields(), qt.DeepEquals, Fields{{"duration", int64(2000)}})
}

type errFields string

func (ef errFields) Error() string {
	return string(ef)
}

func (ef errFields) Fields() Fields {
	return Fields{{"reason", "timeout"}}
}
