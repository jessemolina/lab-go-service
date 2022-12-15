package web

import (
	"context"
	"errors"
	"time"
)

// ================================================================
// FUNCTIONS

// Returns the values from the context.
func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return nil, errors.New("Web valuse missing from context.")
	}
	return v, nil
}

// Returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}
	return v.TraceID
}

// Sets the statuc code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return errors.New("Web valuse missing from context.")
	}
	v.StatusCode = statusCode
	return nil
}

// ================================================================
// TYPES

// Represents the type of value to assign to ctxKey.
type ctxKey int

// Key is used to store and retrieve request values.
const key ctxKey = 1

// Values represents the state of each request.
type Values struct {
	TraceID string
	Now time.Time
	StatusCode int
}
