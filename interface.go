package log

import "time"

// Interface represents the API of both Logger and Entry.
type Interface interface {
	WithFields(Fielder) *EntryFields
	WithField(string, any) *EntryFields
	WithDuration(time.Duration) *EntryFields
	WithError(error) *EntryFields
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
}
