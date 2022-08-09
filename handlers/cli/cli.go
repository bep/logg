// Package cli implements a colored text handler suitable for command-line interfaces.
package cli

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/bep/logg"
	"github.com/fatih/color"
	colorable "github.com/mattn/go-colorable"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

var bold = color.New(color.Bold)

// Colors mapping.
var Colors = [...]*color.Color{
	logg.DebugLevel: color.New(color.FgWhite),
	logg.InfoLevel:  color.New(color.FgBlue),
	logg.WarnLevel:  color.New(color.FgYellow),
	logg.ErrorLevel: color.New(color.FgRed),
}

// Strings mapping.
var Strings = [...]string{
	logg.DebugLevel: "•",
	logg.InfoLevel:  "•",
	logg.WarnLevel:  "•",
	logg.ErrorLevel: "⨯",
}

// Handler implementation.
type Handler struct {
	mu      sync.Mutex
	Writer  io.Writer
	Padding int
}

// New handler.
func New(w io.Writer) *Handler {
	if f, ok := w.(*os.File); ok {
		return &Handler{
			Writer:  colorable.NewColorable(f),
			Padding: 3,
		}
	}

	return &Handler{
		Writer:  w,
		Padding: 3,
	}
}

// HandleLog implements logg.Handler.
func (h *Handler) HandleLog(e *logg.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]

	h.mu.Lock()
	defer h.mu.Unlock()

	color.Fprintf(h.Writer, "%s %-25s", bold.Sprintf("%*s", h.Padding+1, level), e.Message)

	for _, field := range e.Fields {
		if field.Name == "source" {
			continue
		}
		fmt.Fprintf(h.Writer, " %s=%v", color.Sprint(field.Name), field.Value)
	}

	fmt.Fprintln(h.Writer)

	return nil
}
