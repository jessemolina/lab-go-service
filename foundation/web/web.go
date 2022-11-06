package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

// ================================================================
// FUNCTIONS

// Create web app with default mux and shutdown.
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown: shutdown,
	}
}

// ================================================================
// TYPES

// A custom handler function that adds context to http requests.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Web App that provides mux and defaults.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}


// Set a handler function for given HTTP method and path to the application server mux.
func (a *App) Handle(method string, group string, path string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {

		if err := handler(r.Context(), w, r); err != nil {
			// TODO Web app hanlder error handling.
			return
		}

	}

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
