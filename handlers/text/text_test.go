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
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: text.New(&buf, text.Options{Separator: "|"})})
	info := l.WithLevel(logg.LevelInfo)

	info.WithField("user", "tj").WithField("id", "123").Log(logg.String("hello"))
	info.WithField("user", "tj").Log(logg.String("world"))
	info.WithField("user", "tj").WithLevel(logg.LevelError).Log(logg.String("boom"))

	expected := "INFO|hello|user=tj|id=123\nINFO|world|user=tj\nERROR|boom|user=tj\n"

	qt.Assert(t, buf.String(), qt.Equals, expected)
}
