package multi_test

import (
	"testing"

	"github.com/bep/log"
	"github.com/bep/log/handlers/memory"
	"github.com/bep/log/handlers/multi"
	qt "github.com/frankban/quicktest"
)

func TestMulti(t *testing.T) {
	a := memory.New()
	b := memory.New()

	l := log.NewLogger(log.LoggerConfig{
		Level:   log.InfoLevel,
		Handler: multi.New(a, b),
	})

	info := l.WithLevel(log.InfoLevel)

	info.WithField("user", "tj").WithField("id", "123").Log(log.String("hello"))
	info.Log(log.String("world"))
	info.WithLevel(log.ErrorLevel).Log(log.String("boom"))

	qt.Assert(t, a.Entries, qt.HasLen, 3)
	qt.Assert(t, b.Entries, qt.HasLen, 3)
}
