package log_test

import (
	"context"
	"testing"

	"github.com/bep/log"
	qt "github.com/frankban/quicktest"
)

func TestFromContext(t *testing.T) {
	ctx := context.Background()

	logger := log.FromContext(ctx)
	qt.Assert(t, logger, qt.Equals, log.Log)

	logs := log.WithField("foo", "bar")
	ctx = log.NewContext(ctx, logs)

	logger = log.FromContext(ctx)
	qt.Assert(t, logger, qt.Equals, logs)
}
