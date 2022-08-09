package multi_test

import (
	"testing"

	"github.com/bep/logg"
	"github.com/bep/logg/handlers/memory"
	"github.com/bep/logg/handlers/multi"
	qt "github.com/frankban/quicktest"
)

func TestMulti(t *testing.T) {
	a := memory.New()
	b := memory.New()

	l := logg.NewLogger(logg.LoggerConfig{
		Level:   logg.InfoLevel,
		Handler: multi.New(a, b),
	})

	info := l.WithLevel(logg.InfoLevel)

	info.WithField("user", "tj").WithField("id", "123").Log(logg.String("hello"))
	info.Log(logg.String("world"))
	info.WithLevel(logg.ErrorLevel).Log(logg.String("boom"))

	qt.Assert(t, a.Entries, qt.HasLen, 3)
	qt.Assert(t, b.Entries, qt.HasLen, 3)
}
