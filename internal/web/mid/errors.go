package mid

import (
	"context"
	"net/http"

	"github.com/jessemolina/lab-go-service/pkg/web"
	"go.uber.org/zap"
)

// Handles errors coming from outside of the call chains.
// Detects normal application errors which are used to responsd in uniformaty.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.SugaredLogger) web.Middleware {

	// The middleware function that will be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached to the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return nil
		}

		return h
	}

	return m
}
