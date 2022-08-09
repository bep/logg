// Package text implements a development-friendly textual handler.
package text

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/bep/logg"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

// start time.
var start = time.Now()

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

// Colors mapping.
var Colors = [...]int{
	logg.DebugLevel: gray,
	logg.InfoLevel:  blue,
	logg.WarnLevel:  yellow,
	logg.ErrorLevel: red,
}

// Strings mapping.
var Strings = [...]string{
	logg.DebugLevel: "DEBUG",
	logg.InfoLevel:  "INFO",
	logg.WarnLevel:  "WARN",
	logg.ErrorLevel: "ERROR",
}

// Handler implementation.
type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer: w,
	}
}

// HandleLog implements logg.Handler.
func (h *Handler) HandleLog(e *logg.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]

	h.mu.Lock()
	defer h.mu.Unlock()

	ts := time.Since(start) / time.Second
	fmt.Fprintf(h.Writer, "\033[%dm%6s\033[0m[%04d] %-25s", color, level, ts, e.Message)

	for _, f := range e.Fields {
		fmt.Fprintf(h.Writer, " \033[%dm%s\033[0m=%v", color, f.Name, f.Value)
	}

	fmt.Fprintln(h.Writer)

	return nil
}
