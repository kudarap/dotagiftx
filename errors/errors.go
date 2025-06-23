package errors

import (
	"errors"
	"fmt"

	"github.com/kudarap/dotagiftx"
)

// Errors represents application's errors.
type Errors struct {
	Type  dotagiftx.Errors
	Err   error
	Fatal bool
}

// Implements error interface.
func (e *Errors) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Err)
}

// IsEqual checks errors with the same Error Type.
func (e *Errors) IsEqual(t dotagiftx.Errors) bool {
	return errors.Is(t, e.Type)
}

func create(t dotagiftx.Errors, e error, f bool) error {
	return &Errors{t, e, f}
}

// New wraps error into an Errors with Type.
func New(t dotagiftx.Errors, e error) error {
	return create(t, e, false)
}

// Fatal creates a fatal flagged error.
func Fatal(t dotagiftx.Errors, e error) error {
	return create(t, e, true)
}

// Parse returns Errors value if available, else returns nil and ok is false.
// When error is a core.Error type will create new error with that type
// to handle them gracefully. Useful when checking errors types on Parse().
func Parse(err error) (e *Errors, ok bool) {
	// Try packaged error assertion.
	e, ok = err.(*Errors)
	if ok {
		return
	}

	// Try core error assertion as type.
	// handles un-packaged error with valid type that
	// can be use to check typed errors.
	t, ok := err.(dotagiftx.Errors)
	if ok {
		// Error with no details.
		return &Errors{t, errors.New(""), false}, true
	}

	// Cant parse the error.
	return nil, false
}
