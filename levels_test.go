package log

import (
	"encoding/json"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		String string
		Level  Level
		Num    int
	}{
		{"debug", DebugLevel, 0},
		{"info", InfoLevel, 1},
		{"warn", WarnLevel, 2},
		{"warning", WarnLevel, 3},
		{"error", ErrorLevel, 4},
		{"fatal", FatalLevel, 5},
	}

	for _, c := range cases {
		t.Run(c.String, func(t *testing.T) {
			l, err := ParseLevel(c.String)
			qt.Assert(t, err, qt.IsNil, qt.Commentf("parse"))
			qt.Assert(t, l, qt.Equals, c.Level)
		})
	}

	t.Run("invalid", func(t *testing.T) {
		l, err := ParseLevel("something")
		qt.Assert(t, err, qt.Equals, ErrInvalidLevel)
		qt.Assert(t, l, qt.Equals, InvalidLevel)
	})
}

func TestLevel_MarshalJSON(t *testing.T) {
	e := Entry{
		Message: "hello",
		Level:   InfoLevel,
	}

	expect := `{"level":"info","timestamp":"0001-01-01T00:00:00Z","message":"hello"}`

	b, err := json.Marshal(e)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, string(b), qt.Equals, expect)
}

func TestLevel_UnmarshalJSON(t *testing.T) {
	s := `{"fields":[],"level":"info","timestamp":"0001-01-01T00:00:00Z","message":"hello"}`
	e := new(Entry)

	err := json.Unmarshal([]byte(s), e)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, e.Level, qt.Equals, InfoLevel)
}
