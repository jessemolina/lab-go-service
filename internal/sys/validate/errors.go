package validate

import (
	"encoding/json"
	"errors"
)

// ================================================================
// GLOBALS

// Error is for when ID is not in a valid form.
var ErrInvalidID = errors.New("ID is not in the proper form.")


// ================================================================
// FUNCTIONS

// Wraps a provided error with an HTTP status code.
// This should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}

// Iterates through all of the wrapped errors until the root
// error value is reached.
func Cause(err error) error {
	root := err
	for {
		if err = errors.Unwrap(root); err == nil {
			return root
		}
		root = err
	}
}

// ================================================================
// TYPES

// The form used for API response from failures in the API.
type ErrorResponse struct {
	Error  string `json:"error"`
	Fields string `json:"fields,omitempty"`
}

// Passes an error during the request through the app with
// web specific context.
type RequestError struct {
	Err    error
	Status int
	Fields error
}

// Implements the error interface.
// Its used as the default message for the wrapped error.
// This will be shown in the services log
func (err *RequestError) Error() string {
	return err.Err.Error()
}

// Used to indicate an error with a specific request field.
type FieldError struct {
	Fields string `json:"fields"`
	Error  string `json:"error"`
}

// Collection of field errors.
type FieldErrors []FieldError

// Implements the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}

	return string(d)
}
