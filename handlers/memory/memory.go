// Package memory implements an in-memory handler useful for testing, as the
// entries can be accessed after writes.
package memory

import (
	"sync"

	"github.com/bep/logg"
)

// Handler implementation.
type Handler struct {
	mu      sync.Mutex
	Entries []*logg.Entry
}

// New handler.
func New() *Handler {
	return &Handler{}
}

// HandleLog implements logg.Handler.
func (h *Handler) HandleLog(e *logg.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Entries = append(h.Entries, e.Clone())
	return nil
}
