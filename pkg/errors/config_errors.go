package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidToken     = errors.New("invalid bot token")
	ErrStorage          = errors.New("storage error")
	ErrUserNotFound     = errors.New("user not found")
	ErrMessageNotSent   = errors.New("failed to send message")
	ErrInvalidRequest   = errors.New("invalid request")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrContextCanceled  = errors.New("context canceled")
)

type WrappedError struct {
	Err     error
	Message string
	Code    string
}

func (e *WrappedError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *WrappedError) Unwrap() error {
	return e.Err
}

func Wrap(err error, message string, code string) error {
	if err == nil {
		return nil
	}
	return &WrappedError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}