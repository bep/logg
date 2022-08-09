package logg

// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services. See the "handlers"
// directory for implementations.
//
// It is left up to Handlers to implement thread-safety.
type Handler interface {
	// HandleLog is invoked for each log event.
	// Note that if the Entry is going to be used after the call to HandleLog returns,
	// it must be cloned with Clone().
	HandleLog(e *Entry) error
}
