package zapit

import (
	"errors"
	"testing"
)

func TestIsMalformedURLError(t *testing.T) {
	e := errors.New("test")
	if IsMalformedURLError(e) {
		t.Error("Expected error to not be a MalformedURLError type")
	}

	e = &MalformedURLError{}
	if !IsMalformedURLError(e) {
		t.Error("Expected error to be a MalformedURLError type")
	}
}
