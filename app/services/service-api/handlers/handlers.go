package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/jessemolina/ultimate-service/app/services/service-api/handlers/debug/checkgrp"
	"github.com/jessemolina/ultimate-service/app/services/service-api/handlers/v1/testgrp"
	"github.com/jessemolina/ultimate-service/foundation/web"
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
		Log: log,
	}

	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}


// service api mux
func APIMux(cfg APIMuxConfig) *web.App {
	// create a new mux
	app := web.NewApp(cfg.Shutdown)

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	app .Handle(http.MethodGet, "/v1/test", tgh.Test)

	return app
}

// ================================================================
// TYPES

// config contains all mandatory systems required by handlers
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}
