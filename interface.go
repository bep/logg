package log

import (
	"fmt"
	"time"
)

// Interface represents the API of both Logger and Entry.
type Interface interface {
	WithFields(Fielder) *EntryFields
	WithField(string, any) *EntryFields
	WithDuration(time.Duration) *EntryFields
	WithError(error) *EntryFields
	Debug(fmt.Stringer)
	Info(fmt.Stringer)
	Warn(fmt.Stringer)
	Error(fmt.Stringer)
	Fatal(fmt.Stringer)
}
