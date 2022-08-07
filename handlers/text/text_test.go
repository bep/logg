package text_test

import (
	"bytes"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/bep/log"
	"github.com/bep/log/handlers/text"
)

func init() {
	log.Now = func() time.Time {
		return time.Unix(0, 0)
	}
}

func TestTextHandler(t *testing.T) {
	var buf bytes.Buffer

	log.SetHandler(text.New(&buf))
	log.WithField("user", "tj").WithField("id", "123").Info("hello")
	log.WithField("user", "tj").Info("world")
	log.WithField("user", "tj").Error("boom")

	expected := "\x1b[34m  INFO\x1b[0m[0000] hello                     \x1b[34muser\x1b[0m=tj \x1b[34mid\x1b[0m=123\n\x1b[34m  INFO\x1b[0m[0000] world                     \x1b[34muser\x1b[0m=tj\n\x1b[31m ERROR\x1b[0m[0000] boom                      \x1b[31muser\x1b[0m=tj\n"

	qt.Assert(t, buf.String(), qt.Equals, expected)
}
