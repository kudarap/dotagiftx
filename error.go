package dgx

//go:generate stringer -type=Errors -output=errors_string.go

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
