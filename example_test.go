package logg_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bep/logg"
	"github.com/bep/logg/handlers/text"
)

func Example() {
	var buff bytes.Buffer
	// Create a new logger.
	l := logg.New(
		logg.Options{
			Level:   logg.LevelInfo,
			Handler: text.New(&buff, text.Options{Separator: " "}),
		},
	)
	// Create a new log context.
	infoLogger := l.WithLevel(logg.LevelInfo)

	// Logg some user activity.
	userLogger := infoLogger.WithField("user", "foo").WithField("id", "123")
	userLogger.Log(logg.String("logged in"))
	userLogger.WithField("file", "jokes.txt").Log(logg.String("uploaded"))
	userLogger.WithField("file", "morejokes.txt").Log(logg.String("uploaded"))

	fmt.Print(buff.String())

	// Output:
	// INFO logged in user=foo id=123
	// INFO uploaded user=foo id=123 file=jokes.txt
	// INFO uploaded user=foo id=123 file=morejokes.txt
}

func Example_lazyevaluation() {
	var buff bytes.Buffer
	// Create a new logger.
	l := logg.New(
		logg.Options{
			Level:   logg.LevelError,
			Handler: text.New(&buff, text.Options{Separator: " "}),
		},
	)

	errorLogger := l.WithLevel(logg.LevelError)

	// Info is below the logger's level, so
	// nothing will be printed.
	infoLogger := l.WithLevel(logg.LevelInfo)

	// Simulate a busy loop.
	for i := 0; i < 999; i++ {
		infoLogger.WithFields(
			logg.NewFieldsFunc(
				// This func will never be invoked with the current logger's level.
				func() logg.Fields {
					return logg.Fields{
						{"field", strings.Repeat("x", 99)},
					}

				}),
		)
		infoLogger.Log(logg.NewStringFunc(
			// This func will never be invoked with the current logger's level.
			func() string {
				return "log message: " + strings.Repeat("x", 99)
			},
		))

	}

	errorLogger.Log(logg.String("something failed"))

	fmt.Print(buff.String())

	// Output:
	// ERROR something failed

}
