package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/jessemolina/lab-go-service/cmd/services/service-api/handlers/debug/checkgrp"
	"github.com/jessemolina/lab-go-service/cmd/services/service-api/handlers/v1/testgrp"
	"github.com/jessemolina/lab-go-service/pkg/web"
	"github.com/jessemolina/lab-go-service/internal/web/mid"
	"go.uber.org/zap"
)


// ================================================================
// FUNCTIONS

// registers all debug routes from standard library to new mux
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// register std library debug endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	expvar.Handler()

	return mux
}

// debugmux registers both standard library  routes and our own custom
// debug application routes; bypass DefaultServeMux due to security concerns
func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}

	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// Create a web app with service api mux.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
	)

	v1(app, cfg)

	return app
}

// Binds all of the version 1 routes.
func v1(app *web.App, cfg APIMuxConfig) {
	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, "v1", "/test", tgh.Test)

}

// ================================================================
// TYPES

// Config contains all mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}
