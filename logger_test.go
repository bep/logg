package log_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bep/log"
	"github.com/bep/log/handlers/discard"
	"github.com/bep/log/handlers/memory"
	qt "github.com/frankban/quicktest"
)

func TestLogger_printf(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	l.Info(log.String("logged in Tobi"))

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
}

func TestLogger_levels(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	l.Debug(log.String("uploading"))
	l.Info(log.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
}

func TestLogger_WithFields(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	ctx := l.WithFields(log.Fields{{"file", "sloth.png"}})
	ctx.Debug(log.String("uploading"))
	ctx.Info(log.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
	qt.Assert(t, e.Fields, qt.DeepEquals, log.Fields{{"file", "sloth.png"}})
}

func TestLogger_WithField(t *testing.T) {
	h := memory.New()

	l := &log.Logger{
		Handler: h,
		Level:   log.InfoLevel,
	}

	ctx := l.WithField("file", "sloth.png").WithField("user", "Tobi")
	ctx.Debug(log.String("uploading"))
	ctx.Info(log.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
	qt.Assert(t, e.Fields, qt.DeepEquals, log.Fields{{"file", "sloth.png"}, {"user", "Tobi"}})
}

func TestLogger_HandlerFunc(t *testing.T) {
	h := memory.New()
	f := func(e *log.Entry) error {
		return h.HandleLog(e)
	}

	l := &log.Logger{
		Handler: log.HandlerFunc(f),
		Level:   log.InfoLevel,
	}

	l.Info(log.String("logged in Tobi"))

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
}

func BenchmarkLogger_small(b *testing.B) {
	l := &log.Logger{
		Handler: discard.New(),
		Level:   log.InfoLevel,
	}

	for i := 0; i < b.N; i++ {
		l.Info(log.String("login"))
	}
}

func BenchmarkLogger_medium(b *testing.B) {
	l := &log.Logger{
		Handler: discard.New(),
		Level:   log.InfoLevel,
	}

	for i := 0; i < b.N; i++ {
		l.WithFields(log.Fields{
			{"file", "sloth.png"},
			{"type", "image/png"},
			{"size", 1 << 20},
		}).Info(log.String("upload"))
	}
}

func BenchmarkLogger_large(b *testing.B) {
	l := &log.Logger{
		Handler: discard.New(),
		Level:   log.InfoLevel,
	}

	err := fmt.Errorf("boom")

	for i := 0; i < b.N; i++ {
		l.WithFields(log.Fields{
			{"file", "sloth.png"},
			{"type", "image/png"},
			{"size", 1 << 20},
		}).
			WithFields(log.Fields{
				{"some", "more"},
				{"data", "here"},
				{"whatever", "blah blah"},
				{"more", "stuff"},
				{"context", "such useful"},
				{"much", "fun"},
			}).
			WithError(err).Error(log.String("upload failed"))
	}
}

func BenchmarkLogger_levels(b *testing.B) {
	doWork := func(l log.Interface) {
		for i := 0; i < 10; i++ {
			var fn log.StringFunc = func() string {
				return fmt.Sprintf("loging value %s and %s.", "value1", strings.Repeat("value2", i+1))
			}
			l.Info(fn)
		}
	}

	b.Run("level not met, Logger", func(b *testing.B) {
		l := &log.Logger{
			Handler: discard.New(),
			Level:   log.ErrorLevel,
		}
		for i := 0; i < b.N; i++ {
			doWork(l)
		}
	})

	b.Run("level not met, Entry", func(b *testing.B) {
		l := &log.Logger{
			Handler: discard.New(),
			Level:   log.ErrorLevel,
		}
		entry := l.WithField("file", "sloth.png")
		for i := 0; i < b.N; i++ {
			doWork(entry)
		}
	})

	b.Run("level met", func(b *testing.B) {
		l := &log.Logger{
			Handler: discard.New(),
			Level:   log.InfoLevel,
		}
		for i := 0; i < b.N; i++ {
			for j := 0; j < 10; j++ {
				doWork(l)
			}
		}
	})

	b.Run("level met, Entry", func(b *testing.B) {
		l := &log.Logger{
			Handler: discard.New(),
			Level:   log.InfoLevel,
		}
		entry := l.WithField("file", "sloth.png")
		for i := 0; i < b.N; i++ {
			doWork(entry)
		}
	})

}
