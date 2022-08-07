package level_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/bep/log"
	"github.com/bep/log/handlers/level"
	"github.com/bep/log/handlers/memory"
)

func TestLevel(t *testing.T) {
	h := memory.New()
	l := log.NewLogger(
		log.LoggerConfig{Level: log.ErrorLevel, Handler: level.New(h, log.ErrorLevel)},
	)

	info := l.WithLevel(log.InfoLevel)
	info.Log(log.String("hello"))
	info.Log(log.String("world"))
	info.WithLevel(log.ErrorLevel).Log(log.String("boom"))

	qt.Assert(t, h.Entries, qt.HasLen, 1)
	qt.Assert(t, "boom", qt.Equals, h.Entries[0].Message)
}
