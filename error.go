package zapit

import (
	"fmt"
	"strings"
)

// MalformedURLError is an error indicating that the given URL is malformed.
type MalformedURLError struct {
	URL string
}

// Error returns the string representation of the error.
func (e *MalformedURLError) Error() string {
	return fmt.Sprintf("%s is a malformed URL", e.URL)
}

// IsMalformedURLError returns true if the given error is a MalformedURLError type.
// Otherwise, it returns false.
func IsMalformedURLError(err error) bool {
	e, ok := err.(*MalformedURLError)
	return ok && strings.Contains(fmt.Sprint(e), "is a malformed URL")
}
