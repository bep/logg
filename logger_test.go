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

func TestLogger_Log(t *testing.T) {
	h := memory.New()
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: h})
	a := l.WithLevel(logg.LevelInfo)

	a.Log(logg.String("logged in Tobi"))

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, logg.LevelInfo, qt.Equals, e.Level)
}

func TestLogger_Logf(t *testing.T) {
	h := memory.New()
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: h})
	a := l.WithLevel(logg.LevelInfo)

	a.Logf("logged in %s", "Tobi")

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, logg.LevelInfo, qt.Equals, e.Level)
}

func TestLogger_levels(t *testing.T) {
	h := memory.New()
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: h})

	l.WithLevel(logg.LevelTrace).Log(logg.String("uploading"))
	l.WithLevel(logg.LevelDebug).Log(logg.String("uploading"))
	l.WithLevel(logg.LevelInfo).Log(logg.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, logg.LevelInfo, qt.Equals, e.Level)
}

func TestLogger_WithFields(t *testing.T) {
	h := memory.New()
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: h})

	info := l.WithLevel(logg.LevelInfo).WithFields(logg.Fields{{"file", "sloth.png"}})
	info.WithLevel(logg.LevelDebug).Log(logg.String("uploading"))
	info.Log(logg.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, logg.LevelInfo, qt.Equals, e.Level)
	qt.Assert(t, e.Fields, qt.DeepEquals, logg.Fields{{"file", "sloth.png"}})
}

func TestLogger_WithField(t *testing.T) {
	h := memory.New()
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: h})

	info := l.WithLevel(logg.LevelInfo).WithField("file", "sloth.png").WithField("user", "Tobi")
	info.WithLevel(logg.LevelDebug).Log(logg.String("uploading"))
	info.Log(logg.String("upload complete"))

	qt.Assert(t, len(h.Entries), qt.Equals, 1)

	e := h.Entries[0]
	qt.Assert(t, "upload complete", qt.Equals, e.Message)
	qt.Assert(t, logg.LevelInfo, qt.Equals, e.Level)
	qt.Assert(t, e.Fields, qt.DeepEquals, logg.Fields{{"file", "sloth.png"}, {"user", "Tobi"}})
}

func TestLogger_HandlerFunc(t *testing.T) {
	h := memory.New()
	f := func(e *logg.Entry) error {
		return h.HandleLog(e)
	}
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: logg.HandlerFunc(f)})
	info := l.WithLevel(logg.LevelInfo)

	info.Log(logg.String("logged in Tobi"))

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, logg.LevelInfo, qt.Equals, e.Level)
}

func BenchmarkLogger_small(b *testing.B) {
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})
	info := l.WithLevel(logg.LevelInfo)

	for i := 0; i < b.N; i++ {
		info.Log(logg.String("login"))
	}
}

func BenchmarkLogger_medium(b *testing.B) {
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})
	info := l.WithLevel(logg.LevelInfo)

	for i := 0; i < b.N; i++ {
		info.WithFields(logg.Fields{
			{"file", "sloth.png"},
			{"type", "image/png"},
			{"size", 1 << 20},
		}).Log(logg.String("upload"))
	}
}

func BenchmarkLogger_large(b *testing.B) {
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})
	info := l.WithLevel(logg.LevelInfo)

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
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})
	info := l.WithLevel(logg.LevelInfo)
	for i := 0; i < 3; i++ {
		info = info.WithField(fmt.Sprintf("context%d", i), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		info.Log(logg.String("upload"))
	}
}

func BenchmarkLogger_common_context_many_fields(b *testing.B) {
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})
	info := l.WithLevel(logg.LevelInfo)
	for i := 0; i < 42; i++ {
		info = info.WithField(fmt.Sprintf("context%d", i), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		info.Log(logg.String("upload"))
	}
}

func BenchmarkLogger_context_many_fields_duplicate_names_with_field(b *testing.B) {
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		info := l.WithLevel(logg.LevelInfo)
		for i := 0; i < 9999; i++ {
			info = info.WithField("name", "value")
		}
		info.Log(logg.String("upload"))
	}
}

func BenchmarkLogger_context_many_fields_duplicate_names_with_fields(b *testing.B) {
	l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		info := l.WithLevel(logg.LevelInfo)
		for i := 0; i < 3333; i++ {
			info = info.WithFields(logg.Fields{{"name", "value"}, {"name", "value"}, {"name", "value"}})
		}
		info.Log(logg.String("upload"))
	}
}

func BenchmarkLogger_levels(b *testing.B) {
	doWork := func(l logg.LevelLogger) {
		for i := 0; i < 10; i++ {
			l.Log(logg.StringFunc(
				func() string {
					return fmt.Sprintf("loging value %s and %s.", "value1", strings.Repeat("value2", i+1))
				},
			))
		}
	}

	b.Run("level not met", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.New(logg.Options{Level: logg.LevelError, Handler: handlers.Discard})
			error := l.WithLevel(logg.LevelInfo)
			doWork(error)
		}
	})

	b.Run("level not met, one field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.New(logg.Options{Level: logg.LevelError, Handler: handlers.Discard})
			info := l.WithLevel(logg.LevelInfo)
			info = info.WithField("file", "sloth.png")
			doWork(info)
		}
	})

	b.Run("level not met, many fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.New(logg.Options{Level: logg.LevelError, Handler: handlers.Discard})
			info := l.WithLevel(logg.LevelInfo)
			info = info.WithField("file", "sloth.png")
			for i := 0; i < 32; i++ {
				info = info.WithField(fmt.Sprintf("field%d", i), "value")
			}
			doWork(info)
		}
	})

	b.Run("level met", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})
			info := l.WithLevel(logg.LevelInfo)
			for j := 0; j < 10; j++ {
				doWork(info)
			}
		}
	})

	b.Run("level met, one field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})
			info := l.WithLevel(logg.LevelInfo)
			info = info.WithField("file", "sloth.png")
			doWork(info)
		}
	})

	b.Run("level met, many fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l := logg.New(logg.Options{Level: logg.LevelInfo, Handler: handlers.Discard})
			info := l.WithLevel(logg.LevelInfo)
			info = info.WithField("file", "sloth.png")
			for i := 0; i < 32; i++ {
				info = info.WithField(fmt.Sprintf("field%d", i), "value")
			}
			doWork(info)
		}
	})
}
