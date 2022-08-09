// Package text implements a development-friendly textual handler.
package text

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bep/logg"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr, Options{})

// Handler implementation.
type Handler struct {
	opts Options
	w    io.Writer
}

// Options holds options for the text handler.
type Options struct {
	// Separator is the separator between fields.
	// Default is " ".
	Separator string
}

// New handler.
func New(w io.Writer, opts Options) *Handler {
	if opts.Separator == "" {
		opts.Separator = " "
	}
	return &Handler{
		w:    w,
		opts: opts,
	}
}

// HandleLog implements logg.Handler.
func (h *Handler) HandleLog(e *logg.Entry) error {
	fields := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		fields[i] = fmt.Sprintf("%s=%v", f.Name, f.Value)
	}

	fmt.Fprintf(h.w, "%s%s%s%s%s\n", strings.ToUpper(e.Level.String()), h.opts.Separator, e.Message, h.opts.Separator, strings.Join(fields, h.opts.Separator))

	return nil
}
