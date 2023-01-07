// Package web provides a small web framework extension.
package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

// ================================================================
// FUNCTIONS

// Create a web application with embedded mux to handle routes and middlware.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// ================================================================
// TYPES

// A handler function type used to handle http request with context.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// A web application type that embeds a mux, signal, and middleware.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// Sets a handler function for given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {

	// Wrap specific middleware around the provided handler.
	handler = wrapMiddleware(mw, handler)

	// Wrap the application's general middeware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	// The function for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx = context.WithValue(ctx, key, &v)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}

	}

	// Set the endpoint's full path.
	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}

	// Handle
	a.ContextMux.Handle(method, finalPath, h)

}

// Shutdown web application.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
