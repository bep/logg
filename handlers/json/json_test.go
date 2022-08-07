package json_test

import (
	"bytes"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/bep/clocks"
	"github.com/bep/log"
	"github.com/bep/log/handlers/json"
)

func TestJSONHandler(t *testing.T) {
	var buf bytes.Buffer

	l := log.NewLogger(
		log.LoggerConfig{
			Level:   log.InfoLevel,
			Handler: json.New(&buf),
			Clock:   clocks.Fixed(clocks.TimeCupFinalNorway1976),
		})

	info := l.WithLevel(log.InfoLevel)

	info.WithField("user", "tj").WithField("id", "123").Log(log.String("hello"))
	info.Log(log.String("world"))
	info.WithLevel(log.ErrorLevel).Log(log.String("boom"))

	expected := "{\"fields\":{\"id\":\"123\",\"user\":\"tj\"},\"level\":\"info\",\"timestamp\":\"1976-10-24T12:15:02.127686412Z\",\"message\":\"hello\"}\n{\"fields\":{},\"level\":\"info\",\"timestamp\":\"1976-10-24T12:15:02.127686412Z\",\"message\":\"world\"}\n{\"fields\":{},\"level\":\"error\",\"timestamp\":\"1976-10-24T12:15:02.127686412Z\",\"message\":\"boom\"}\n"

	qt.Assert(t, buf.String(), qt.Equals, expected)
}
