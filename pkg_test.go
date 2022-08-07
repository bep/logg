package log_test

import (
	"errors"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/bep/log"
	"github.com/bep/log/handlers/memory"
)

type Pet struct {
	Name string
	Age  int
}

func (p *Pet) Fields() log.Fields {
	return log.Fields{
		log.Field{
			"name", p.Name,
		},
		log.Field{
			"age", p.Age,
		},
	}
}

func TestInfo(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	log.Infof("logged in %s", "Tobi")

	e := h.Entries[0]
	qt.Assert(t, "logged in Tobi", qt.Equals, e.Message)
	qt.Assert(t, log.InfoLevel, qt.Equals, e.Level)
}

func TestFielder(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	pet := &Pet{"Tobi", 3}
	log.WithFields(pet).Info("add pet")

	e := h.Entries[0]
	qt.Assert(t, e.Fields, qt.DeepEquals, log.Fields{
		{"name", "Tobi"},
		{"age", 3},
	})
}

// Unstructured logging is supported, but not recommended since it is hard to query.
func Example_unstructured() {
	log.Infof("%s logged in", "Tobi")
}

// Structured logging is supported with fields, and is recommended over the formatted message variants.
func Example_structured() {
	log.WithField("user", "Tobo").Info("logged in")
}

// Errors are passed to WithError(), populating the "error" field.
func Example_errors() {
	err := errors.New("boom")
	log.WithError(err).Error("upload failed")
}

// Multiple fields can be set, via chaining, or WithFields().
func Example_multipleFields() {
	log.WithFields(log.Fields{
		{"user", "Tobi"},
		{"file", "sloth.png"},
		{"type", "image/png"},
	}).Info("upload")
}
