package kvdb

import (
	"errors"
	"fmt"
)

// Store defines the kvdb public interface
type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Flush() error
	Close() error
	IsNotFoundError(err error) bool
	IsBadRequestError(err error) bool
}

// NotFoundError indicates that no value was found for the given key
type NotFoundError struct {
	missingKey string
}

// NewNotFoundError returns a new error for the missing key
func NewNotFoundError(missingKey string) error {
	return &NotFoundError{missingKey}
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("Could not find value for key: %s", n.missingKey)
}

// IsNotFoundError returns true if the error, or any of the wrapped errors
// is of type BadRequestError
func IsNotFoundError(err error) bool {
	var notFoundError *NotFoundError
	return errors.As(err, &notFoundError)
}

// BadRequestError represents an error by the consumer of the database
type BadRequestError struct {
	message string
}

// NewBadRequestError returns a new BadRequestError given an error message
func NewBadRequestError(message string) error {
	return &BadRequestError{message}
}

func (b *BadRequestError) Error() string {
	return b.message
}

// IsBadRequestError returns true if the error, or any of the wrapped errors
// is of type BadRequestError
func IsBadRequestError(err error) bool {
	var badRequestError *BadRequestError
	return errors.As(err, &badRequestError)
}
