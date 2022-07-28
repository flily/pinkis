package meta

import (
	"reflect"
	"testing"

	"errors"
)

func TestMetaError(t *testing.T) {
	err := NewMetaError("hello, world")

	message := "hello, world"
	if err.Error() != message {
		t.Errorf("error message: %s", err)
	}

	if !errors.Is(err, ErrMetaError) {
		t.Errorf("err (%v) is not ErrMetaError (%v)", err, ErrMetaError)
	}
}

func TestNotDuplicatableError(t *testing.T) {
	err := NewNotDuplicatableError(reflect.Chan)

	message := "kind chan is not duplicatable"
	if err.Error() != message {
		t.Errorf("error message: %s", err)
	}
}
