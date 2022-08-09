package level_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/bep/logg"
	"github.com/bep/logg/handlers/level"
	"github.com/bep/logg/handlers/memory"
)

func TestLevel(t *testing.T) {
	h := memory.New()
	l := logg.NewLogger(
		logg.LoggerConfig{Level: logg.LevelError, Handler: level.New(h, logg.LevelError)},
	)

	info := l.WithLevel(logg.LevelInfo)
	info.Log(logg.String("hello"))
	info.Log(logg.String("world"))
	info.WithLevel(logg.LevelError).Log(logg.String("boom"))

	qt.Assert(t, h.Entries, qt.HasLen, 1)
	qt.Assert(t, "boom", qt.Equals, h.Entries[0].Message)
}
