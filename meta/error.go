package meta

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrMetaError      = errors.New("meta error")
	ErrNotDuplcatable = NewMetaError("not duplicatable")
	ErrUntypedNil  = NewMetaError("untyped nil is unacceptable")
)

type MetaError struct {
	Base    error
	Message string
}

func (e MetaError) Error() string {
	return e.Message
}

func (e MetaError) Unwrap() error {
	return e.Base
}

func NewMetaError(format string, args ...interface{}) error {
	return MetaError{
		Base:    ErrMetaError,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewNotDuplicatableError(kind reflect.Kind) error {
	return &MetaError{
		Base:    ErrNotDuplcatable,
		Message: fmt.Sprintf("kind %s is not duplicatable", kind),
	}
}
