package web

import (
	"errors"
)


// ================================================================
// FUNCTIONS

// Returns an error that causes the framework to signal
// a graceful shutdown.
func NewShutdownError(message string) error {
	return &shutdownError{message}

}

// Checks to see if the shutdown error is contained
// in the specified error value.
func IsShutdown(err error) bool {
	var se *shutdownError
	return errors.As(err, &se)
}

// ================================================================
// TYPES

// Type used to help with the graceful termination of the service.
type shutdownError struct {
	Message string
}

// The implementation of the error interface.
func (se *shutdownError) Error() string {
	return se.Message
}
