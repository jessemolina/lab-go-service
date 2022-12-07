package checkgrp

import (
	"encoding/json"
	"net/http"
	"os"

	"go.uber.org/zap"
)

// ================================================================
// FUNCTIONS

// converts data struct to http response with proper header and status code
func response(w http.ResponseWriter, statusCode int, data interface{}) error {

	// convert the response to JSON
	jsonData, err := json.Marshal(data)
	if err != nil{
		return err
	}

	// set the content type and headers after successful marshaling
	w.Header().Set("Content-Type", "application/json")

	// write the status code to the respone
	w.WriteHeader(statusCode)

	// send the result back to the client
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

// ================================================================
// TYPES

type Handlers struct {
	Build string
	Log *zap.SugaredLogger
}

// checks that the application services are ready
func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}
	statusCode := http.StatusOK

	if err := response(w, statusCode, data); err != nil{
		h.Log.Errorw("readiness", "ERROR", err)
	}

	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

}


// checks that the service is alive
func (h Handlers) Liveness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     h.Build,
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),

	}

	statusCode := http.StatusOK

	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("liveness", "ERROR", err)
	}

	h.Log.Infow("liveness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

}
