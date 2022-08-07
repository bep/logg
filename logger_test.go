package log_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bep/log"
	"github.com/bep/log/handlers/memory"
	qt "github.com/frankban/quicktest"
)

func TestLogger_printf(t *testing.T) {
	h := memory.New()
	l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: h})
	a := l.WithLevel(log.InfoLevel)

	a.Log(log.String("logged in Tobi"))

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
}

func TestLogger_levels(t *testing.T) {
	h := memory.New()
	l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: h})

	l.WithLevel(log.DebugLevel).Log(log.String("uploading"))
	l.WithLevel(log.InfoLevel).Log(log.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
}

func TestLogger_WithFields(t *testing.T) {
	h := memory.New()
	l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: h})

	info := l.WithLevel(log.InfoLevel).WithFields(log.Fields{{"file", "sloth.png"}})
	info.WithLevel(log.DebugLevel).Log(log.String("uploading"))
	info.Log(log.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
	qt.Assert(t, e.Fields, qt.DeepEquals, log.Fields{{"file", "sloth.png"}})
}

func TestLogger_WithField(t *testing.T) {
	h := memory.New()
	l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: h})

	info := l.WithLevel(log.InfoLevel).WithField("file", "sloth.png").WithField("user", "Tobi")
	info.WithLevel(log.DebugLevel).Log(log.String("uploading"))
	info.Log(log.String("upload complete"))

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
	l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: log.HandlerFunc(f)})
	info := l.WithLevel(log.InfoLevel)

	info.Log(log.String("logged in Tobi"))

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
}

func BenchmarkLogger_small(b *testing.B) {
	l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: log.NoopHandler})
	info := l.WithLevel(log.InfoLevel)

	for i := 0; i < b.N; i++ {
		info.Log(log.String("login"))
	}
}

func BenchmarkLogger_medium(b *testing.B) {
	l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: log.NoopHandler})
	info := l.WithLevel(log.InfoLevel)

	for i := 0; i < b.N; i++ {
		info.WithFields(log.Fields{
			{"file", "sloth.png"},
			{"type", "image/png"},
			{"size", 1 << 20},
		}).Log(log.String("upload"))
	}
}

func BenchmarkLogger_large(b *testing.B) {
	l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: log.NoopHandler})
	info := l.WithLevel(log.InfoLevel)

	err := fmt.Errorf("boom")

	for i := 0; i < b.N; i++ {
		info.WithFields(log.Fields{
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
			WithError(err).Log(log.String("upload failed"))
	}
}

func BenchmarkLogger_levels(b *testing.B) {
	doWork := func(l log.Logger) {
		for i := 0; i < 10; i++ {
			l.Log(log.NewStringFunc(
				func() string {
					return fmt.Sprintf("loging value %s and %s.", "value1", strings.Repeat("value2", i+1))
				},
			))
		}
	}

	b.Run("level not met", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := log.NewLogger(log.LoggerConfig{Level: log.ErrorLevel, Handler: log.NoopHandler})
			error := l.WithLevel(log.InfoLevel)
			doWork(error)
		}
	})

	b.Run("level not met, one field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := log.NewLogger(log.LoggerConfig{Level: log.ErrorLevel, Handler: log.NoopHandler})
			info := l.WithLevel(log.InfoLevel)
			info = info.WithField("file", "sloth.png")
			doWork(info)
		}
	})

	b.Run("level not met, many fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := log.NewLogger(log.LoggerConfig{Level: log.ErrorLevel, Handler: log.NoopHandler})
			info := l.WithLevel(log.InfoLevel)
			info = info.WithField("file", "sloth.png")
			for i := 0; i < 32; i++ {
				info = info.WithField(fmt.Sprintf("field%d", i), "value")
			}
			doWork(info)
		}
	})

	b.Run("level met", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: log.NoopHandler})
			info := l.WithLevel(log.InfoLevel)
			for j := 0; j < 10; j++ {
				doWork(info)
			}
		}
	})

	b.Run("level met, one field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: log.NoopHandler})
			info := l.WithLevel(log.InfoLevel)
			info = info.WithField("file", "sloth.png")
			doWork(info)
		}
	})

	b.Run("level met, many fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := log.NewLogger(log.LoggerConfig{Level: log.InfoLevel, Handler: log.NoopHandler})
			info := l.WithLevel(log.InfoLevel)
			info = info.WithField("file", "sloth.png")
			for i := 0; i < 32; i++ {
				info = info.WithField(fmt.Sprintf("field%d", i), "value")
			}
			doWork(info)
		}
	})
}
