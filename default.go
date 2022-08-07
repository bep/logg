package log

import (
	"bytes"
	"fmt"
	"log"
)

// handleStdLog outpouts to the stlib log.
func handleStdLog(e *Entry) error {
	level := levelNames[e.Level]

	var b bytes.Buffer
	fmt.Fprintf(&b, "%5s %-25s", level, e.Message)

	for _, f := range e.Fields {
		fmt.Fprintf(&b, " %s=%v", f.Name, f.Value)
	}

	log.Println(b.String())

	return nil
}
