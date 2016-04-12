package vcs

import (
	"errors"
	"testing"
)

func TestNewGetError(t *testing.T) {
	base := errors.New("Foo error")
	out := "This is a test"

	e := NewGetError(base, out)

	switch e.(type) {
	case *GetError:
		// This is the right error type
	default:
		t.Error("Wrong error type returned from NewGetError")
	}
}
