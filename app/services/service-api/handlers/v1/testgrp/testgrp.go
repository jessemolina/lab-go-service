package testgrp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// ================================================================
// TYPES

type Handlers struct {
	Build string
	Log *zap.SugaredLogger
}


func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(status)

	statusCode := http.StatusOK
	h.Log.Infow("test", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}
