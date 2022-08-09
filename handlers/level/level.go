// Package level implements a level filter handler.
package level

import "github.com/bep/logg"

// Handler implementation.
type Handler struct {
	Level   logg.Level
	Handler logg.Handler
}

// New handler.
func New(h logg.Handler, level logg.Level) *Handler {
	return &Handler{
		Level:   level,
		Handler: h,
	}
}

// HandleLog implements logg.Handler.
func (h *Handler) HandleLog(e *logg.Entry) error {
	if e.Level < h.Level {
		return nil
	}

	return h.Handler.HandleLog(e)
}
