package dbErrors

import (
	gerrors "github.com/pkg/errors"
)
const (
	ErrDoesNotExists = "Resource does not exists"
	ErrAlreadyExists = "Resource already exists"
	ErrorUnknown =     "Unknown database error"
	ErrorViolates =    "Constraint violating"
)

type DbErr interface {
	Error() string
	Code() string
}

type DatabaseError struct {
	ErrError error
	ErrCode string
}

func (d DatabaseError) Code() string {
	return d.ErrCode
}

func (d DatabaseError) Error() string {
	return d.ErrError.Error()
}

func NewDoesNotExists(err error, msg string) error {
	return DatabaseError{
		ErrError: gerrors.Wrap(err, msg),
		ErrCode:  ErrDoesNotExists,
	}
}

func NewUnknown(err error, msg string) error {
	return DatabaseError{
		ErrError: gerrors.Wrap(err, msg),
		ErrCode:  ErrorUnknown,
	}
}

func NewAlreadyExists(err error, msg string) error {
	return DatabaseError{
		ErrError: gerrors.Wrap(err, msg),
		ErrCode:  ErrAlreadyExists,
	}
}

func NewViolates(err error, msg string) error {
	return DatabaseError{
		ErrError: gerrors.Wrap(err, msg),
		ErrCode:  ErrorViolates,
	}
}