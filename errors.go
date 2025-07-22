//go:generate go tool stringer -type=Errors -output=errors_string.go

package dotagiftx

import (
	"fmt"
)

// Error indexes are used for auto-increment identifier for error code generation.
// The enumeration below is to avoid conflict.
const (
	storageErrorIndex   = 100
	authErrorIndex      = 1000
	userErrorIndex      = 1100
	marketErrorIndex    = 2100
	catalogErrorIndex   = 2200
	itemErrorIndex      = 3000
	trackErrorIndex     = 4000
	reportErrorIndex    = 5000
	deliveryErrorIndex  = 6000
	inventoryErrorIndex = 6100
)

var appErrorText = map[Errors]string{}

// Errors represents app's error.
type Errors uint

// Error implements error interface.
func (i Errors) Error() string {
	return appErrorText[i]
}

// Code returns error code.
func (i Errors) Code() string {
	return i.String()
}

func (i Errors) X(err error) XErrors {
	return XErrors{Type: i, Err: err}
}

// XErrors represents application's errors.
type XErrors struct {
	Type  Errors
	Err   error
	Fatal bool
}

// Implements error interface.
func (x XErrors) Error() string {
	return fmt.Sprintf("%s: %s", x.Type, x.Err)
}

func NewXError(t Errors, err error) XErrors {
	return XErrors{Type: t, Err: err}
}
