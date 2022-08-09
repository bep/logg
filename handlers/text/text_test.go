package text_test

import (
	"bytes"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/bep/logg"
	"github.com/bep/logg/handlers/text"
)

func TestTextHandler(t *testing.T) {
	var buf bytes.Buffer
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: text.New(&buf)})
	info := l.WithLevel(logg.LevelInfo)

	info.WithField("user", "tj").WithField("id", "123").Log(logg.String("hello"))
	info.WithField("user", "tj").Log(logg.String("world"))
	info.WithField("user", "tj").WithLevel(logg.LevelError).Log(logg.String("boom"))

	expected := "\x1b[34m  INFO\x1b[0m[0000] hello                     \x1b[34muser\x1b[0m=tj \x1b[34mid\x1b[0m=123\n\x1b[34m  INFO\x1b[0m[0000] world                     \x1b[34muser\x1b[0m=tj\n\x1b[31m ERROR\x1b[0m[0000] boom                      \x1b[31muser\x1b[0m=tj\n"

	qt.Assert(t, buf.String(), qt.Equals, expected)
}
