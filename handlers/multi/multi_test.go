package multi_test

import (
	"testing"
	"time"

	"github.com/bep/log"
	"github.com/bep/log/handlers/memory"
	"github.com/bep/log/handlers/multi"
	qt "github.com/frankban/quicktest"
)

func init() {
	log.Now = func() time.Time {
		return time.Unix(0, 0)
	}
}

func Test(t *testing.T) {
	a := memory.New()
	b := memory.New()

	log.SetHandler(multi.New(a, b))
	log.WithField("user", "tj").WithField("id", "123").Info(log.String("hello"))
	log.Info(log.String("world"))
	log.Error(log.String("boom"))

	qt.Assert(t, a.Entries, qt.HasLen, 3)
	qt.Assert(t, b.Entries, qt.HasLen, 3)
}
