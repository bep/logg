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

	l := logg.New(logg.Options{
		Level:   logg.LevelInfo,
		Handler: multi.New(a, b),
	})

	info := l.WithLevel(logg.LevelInfo)

	info.WithField("user", "tj").WithField("id", "123").Log(logg.String("hello"))
	info.Log(logg.String("world"))
	info.WithLevel(logg.LevelError).Log(logg.String("boom"))

	qt.Assert(t, a.Entries, qt.HasLen, 3)
	qt.Assert(t, b.Entries, qt.HasLen, 3)
}

func TestMultiModifyEntry(t *testing.T) {
	var a logg.HandlerFunc = func(e *logg.Entry) error {
		e.Message += "-modified"
		e.Fields = append(e.Fields, logg.Field{Name: "added", Value: "value"})
		return nil
	}

	b := memory.New()

	l := logg.New(
		logg.Options{
			Level:   logg.LevelInfo,
			Handler: multi.New(a, b),
		})

	l.WithLevel(logg.LevelInfo).WithField("initial", "value").Log(logg.String("text"))

	qt.Assert(t, b.Entries, qt.HasLen, 1)
	qt.Assert(t, b.Entries[0].Message, qt.Equals, "text-modified")
	qt.Assert(t, b.Entries[0].Fields, qt.HasLen, 2)
	qt.Assert(t, b.Entries[0].Fields[0].Name, qt.Equals, "initial")
	qt.Assert(t, b.Entries[0].Fields[1].Name, qt.Equals, "added")
}

func TestMultDisableEntry(t *testing.T) {
	var a logg.HandlerFunc = func(e *logg.Entry) error {
		if e.Fields[0].Value == "v2" {
			e.Disabled = true
		}
		return nil
	}

	b := memory.New()

	l := logg.New(
		logg.Options{
			Level:   logg.LevelInfo,
			Handler: multi.New(a, b),
		})

	l.WithLevel(logg.LevelInfo).WithField("v", "v1").Log(logg.String("text1"))
	l.WithLevel(logg.LevelInfo).WithField("v", "v2").Log(logg.String("text2"))
	l.WithLevel(logg.LevelInfo).WithField("v", "v3").Log(logg.String("text3"))

	qt.Assert(t, b.Entries, qt.HasLen, 2)
}
