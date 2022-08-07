package level_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/bep/log"
	"github.com/bep/log/handlers/level"
	"github.com/bep/log/handlers/memory"
)

func Test(t *testing.T) {
	h := memory.New()

	ctx := log.Logger{
		Handler: level.New(h, log.ErrorLevel),
		Level:   log.InfoLevel,
	}

	ctx.Info(log.String("hello"))
	ctx.Info(log.String("world"))
	ctx.Error(log.String("boom"))

	qt.Assert(t, h.Entries, qt.HasLen, 1)
	qt.Assert(t, "boom", qt.Equals, h.Entries[0].Message)
}
