package web

// ================================================================
// TYPES

// Function designed to run code before and/or after another Handler.
type Middleware func(Handler) Handler

// ================================================================
// FUNCTIONS

// Creates a new handler that wraps functions around the final Handler.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {

	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
