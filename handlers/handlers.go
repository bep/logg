package handlers

import "github.com/bep/logg"

// Discard is a no-op handler that discards all log messages.
var Discard = logg.HandlerFunc(func(e *logg.Entry) error {
	return nil
})
