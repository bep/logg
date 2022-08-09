package logg_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bep/logg"
	"github.com/bep/logg/handlers"
	"github.com/bep/logg/handlers/memory"
	qt "github.com/frankban/quicktest"
)

func TestLogger_printf(t *testing.T) {
	h := memory.New()
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: h})
	a := l.WithLevel(logg.InfoLevel)

	a.Log(logg.String("logged in Tobi"))

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, logg.InfoLevel, qt.Equals, e.Level)
}

func TestLogger_levels(t *testing.T) {
	h := memory.New()
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: h})

	l.WithLevel(logg.DebugLevel).Log(logg.String("uploading"))
	l.WithLevel(logg.InfoLevel).Log(logg.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, logg.InfoLevel, qt.Equals, e.Level)
}

func TestLogger_WithFields(t *testing.T) {
	h := memory.New()
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: h})

	info := l.WithLevel(logg.InfoLevel).WithFields(logg.Fields{{"file", "sloth.png"}})
	info.WithLevel(logg.DebugLevel).Log(logg.String("uploading"))
	info.Log(logg.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, logg.InfoLevel, qt.Equals, e.Level)
	qt.Assert(t, e.Fields, qt.DeepEquals, logg.Fields{{"file", "sloth.png"}})
}

func TestLogger_WithField(t *testing.T) {
	h := memory.New()
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: h})

	info := l.WithLevel(logg.InfoLevel).WithField("file", "sloth.png").WithField("user", "Tobi")
	info.WithLevel(logg.DebugLevel).Log(logg.String("uploading"))
	info.Log(logg.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, logg.InfoLevel, qt.Equals, e.Level)
	qt.Assert(t, e.Fields, qt.DeepEquals, logg.Fields{{"file", "sloth.png"}, {"user", "Tobi"}})
}

func TestLogger_HandlerFunc(t *testing.T) {
	h := memory.New()
	f := func(e *logg.Entry) error {
		return h.HandleLog(e)
	}
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: logg.HandlerFunc(f)})
	info := l.WithLevel(logg.InfoLevel)

	info.Log(logg.String("logged in Tobi"))

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, logg.InfoLevel, qt.Equals, e.Level)
}

func BenchmarkLogger_small(b *testing.B) {
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: handlers.Discard})
	info := l.WithLevel(logg.InfoLevel)

	for i := 0; i < b.N; i++ {
		info.Log(logg.String("login"))
	}
}

func BenchmarkLogger_medium(b *testing.B) {
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: handlers.Discard})
	info := l.WithLevel(logg.InfoLevel)

	for i := 0; i < b.N; i++ {
		info.WithFields(logg.Fields{
			{"file", "sloth.png"},
			{"type", "image/png"},
			{"size", 1 << 20},
		}).Log(logg.String("upload"))
	}
}

func BenchmarkLogger_large(b *testing.B) {
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: handlers.Discard})
	info := l.WithLevel(logg.InfoLevel)

	err := fmt.Errorf("boom")

	for i := 0; i < b.N; i++ {
		info.WithFields(logg.Fields{
			{"file", "sloth.png"},
			{"type", "image/png"},
			{"size", 1 << 20},
		}).
			WithFields(logg.Fields{
				{"some", "more"},
				{"data", "here"},
				{"whatever", "blah blah"},
				{"more", "stuff"},
				{"context", "such useful"},
				{"much", "fun"},
			}).
			WithError(err).Log(logg.String("upload failed"))
	}
}

func BenchmarkLogger_common_context(b *testing.B) {
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: handlers.Discard})
	info := l.WithLevel(logg.InfoLevel)
	for i := 0; i < 3; i++ {
		info = info.WithField(fmt.Sprintf("context%d", i), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		info.Log(logg.String("upload"))
	}
}

func BenchmarkLogger_common_context_many_fields(b *testing.B) {
	l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: handlers.Discard})
	info := l.WithLevel(logg.InfoLevel)
	for i := 0; i < 42; i++ {
		info = info.WithField(fmt.Sprintf("context%d", i), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		info.Log(logg.String("upload"))
	}
}

func BenchmarkLogger_levels(b *testing.B) {
	doWork := func(l logg.LevelLogger) {
		for i := 0; i < 10; i++ {
			l.Log(logg.NewStringFunc(
				func() string {
					return fmt.Sprintf("loging value %s and %s.", "value1", strings.Repeat("value2", i+1))
				},
			))
		}
	}

	b.Run("level not met", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.NewLogger(logg.LoggerConfig{Level: logg.ErrorLevel, Handler: handlers.Discard})
			error := l.WithLevel(logg.InfoLevel)
			doWork(error)
		}
	})

	b.Run("level not met, one field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.NewLogger(logg.LoggerConfig{Level: logg.ErrorLevel, Handler: handlers.Discard})
			info := l.WithLevel(logg.InfoLevel)
			info = info.WithField("file", "sloth.png")
			doWork(info)
		}
	})

	b.Run("level not met, many fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.NewLogger(logg.LoggerConfig{Level: logg.ErrorLevel, Handler: handlers.Discard})
			info := l.WithLevel(logg.InfoLevel)
			info = info.WithField("file", "sloth.png")
			for i := 0; i < 32; i++ {
				info = info.WithField(fmt.Sprintf("field%d", i), "value")
			}
			doWork(info)
		}
	})

	b.Run("level met", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: handlers.Discard})
			info := l.WithLevel(logg.InfoLevel)
			for j := 0; j < 10; j++ {
				doWork(info)
			}
		}
	})

	b.Run("level met, one field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: handlers.Discard})
			info := l.WithLevel(logg.InfoLevel)
			info = info.WithField("file", "sloth.png")
			doWork(info)
		}
	})

	b.Run("level met, many fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.NewLogger(logg.LoggerConfig{Level: logg.InfoLevel, Handler: handlers.Discard})
			info := l.WithLevel(logg.InfoLevel)
			info = info.WithField("file", "sloth.png")
			for i := 0; i < 32; i++ {
				info = info.WithField(fmt.Sprintf("field%d", i), "value")
			}
			doWork(info)
		}
	})
}
