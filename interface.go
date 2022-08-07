package log

import (
	"fmt"
	"time"
)

type Leveler interface {
	WithLevel(Level) *EntryFields
}

type Logger interface {
	Log(fmt.Stringer)
	Leveler
	WithFields(Fielder) *EntryFields
	WithField(string, any) *EntryFields
	WithDuration(time.Duration) *EntryFields
	WithError(error) *EntryFields
}
