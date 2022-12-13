package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

// ================================================================
// TYPES

// A custom handler function that adds context to http requests.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Web App that embeds a mux, signals, and middleware.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// Set a handler function for given HTTP method and path to the application server mux.
func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {

	// Wrap mw specific handlers around the original handler.
	handler = wrapMiddleware(mw, handler)

	// Wrap the application's general middeware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	// The function for each request.
	h := func(w http.ResponseWriter, r *http.Request) {

		if err := handler(r.Context(), w, r); err != nil {
			return
		}

	}

	// Set the endpoint's full path.
	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}

	a.ContextMux.Handle(method, finalPath, h)

}

// Shutdown web application.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// ================================================================
// FUNCTIONS

// Create web app with default mux and shutdown.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}
