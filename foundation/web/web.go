package web

import (
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

// ================================================================
// FUNCTIONS

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown: shutdown,
	}
}

// ================================================================
// TYPES

type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
