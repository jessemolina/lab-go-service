package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
)

// registers all debug routes from standard library
// using a new serve mux
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