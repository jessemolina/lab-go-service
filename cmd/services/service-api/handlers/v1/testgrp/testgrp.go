package testgrp

import (
	"context"
	"net/http"

	"github.com/jessemolina/ultimate-service/foundation/web"
	"go.uber.org/zap"
)

// ================================================================
// TYPES

// Handlers manages the set of test endpoints.
type Handlers struct {
	Build string
	Log *zap.SugaredLogger
}

// Test handler used for development.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	statusCode := http.StatusOK
	h.Log.Infow("test", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	return web.Respond(ctx, w, status, statusCode)
}
