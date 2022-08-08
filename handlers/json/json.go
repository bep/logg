// Package json implements a JSON handler.
package json

import (
	"encoding/json"
	"io"

	"github.com/bep/log"
)

type Handler struct {
	w io.Writer
}

// New Handler implementation for JSON logging.
// Eeach log Entry is written as a single JSON object, no more than one write to w.
// The writer w should be safe for concurrent use by multiple
// goroutines if the returned Handler will be used concurrently.
func New(w io.Writer) *Handler {
	return &Handler{
		w,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	enc := json.NewEncoder(h.w)
	enc.SetEscapeHTML(false)
	return enc.Encode(e)
}
