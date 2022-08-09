// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package benchmarks

import (
	"io"

	"github.com/bep/logg"
	"github.com/bep/logg/handlers/json"
)

func newDisabledLoggLog() logg.LevelLogger {
	logger := logg.NewLogger(logg.LoggerConfig{
		Handler: json.New(io.Discard),
		Level:   logg.ErrorLevel,
	})
	return logger.WithLevel(logg.InfoLevel)
}

func newLoggLog() logg.LevelLogger {
	logger := logg.NewLogger(logg.LoggerConfig{
		Handler: json.New(io.Discard),
		Level:   logg.DebugLevel,
	})

	return logger.WithLevel(logg.DebugLevel)
}

func fakeLognFields() logg.FieldsFunc {
	return func() logg.Fields {
		return logg.Fields{
			{Name: "int", Value: _tenInts[0]},
			{Name: "ints", Value: _tenInts},
			{Name: "string", Value: _tenStrings[0]},
			{Name: "strings", Value: _tenStrings},
			{Name: "time", Value: _tenTimes[0]},
			{Name: "times", Value: _tenTimes},
			{Name: "user1", Value: _oneUser},
			{Name: "user2", Value: _oneUser},
			{Name: "users", Value: _tenUsers},
			{Name: "error", Value: errExample},
		}
	}
}
