package handlers

import "github.com/bep/log"

// Discard is a no-op handler that discards all log messages.
var Discard = log.HandlerFunc(func(e *log.Entry) error {
	return nil
})
