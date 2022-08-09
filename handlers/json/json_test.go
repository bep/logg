package json_test

import (
	"bytes"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/bep/clocks"
	"github.com/bep/logg"
	"github.com/bep/logg/handlers/json"
)

func TestJSONHandler(t *testing.T) {
	var buf bytes.Buffer

	l := logg.NewLogger(
		logg.LoggerConfig{
			Level:   logg.LevelInfo,
			Handler: json.New(&buf),
			Clock:   clocks.Fixed(clocks.TimeCupFinalNorway1976),
		})

	info := l.WithLevel(logg.LevelInfo)

	info.WithField("user", "tj").WithField("id", "123").Log(logg.String("hello"))
	info.Log(logg.String("world"))
	info.WithLevel(logg.LevelError).Log(logg.String("boom"))

	expected := "{\"level\":\"info\",\"timestamp\":\"1976-10-24T12:15:02.127686412Z\",\"fields\":[{\"name\":\"user\",\"value\":\"tj\"},{\"name\":\"id\",\"value\":\"123\"}],\"message\":\"hello\"}\n{\"level\":\"info\",\"timestamp\":\"1976-10-24T12:15:02.127686412Z\",\"message\":\"world\"}\n{\"level\":\"error\",\"timestamp\":\"1976-10-24T12:15:02.127686412Z\",\"message\":\"boom\"}\n"

	qt.Assert(t, buf.String(), qt.Equals, expected)
}
