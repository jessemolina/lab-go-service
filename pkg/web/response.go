package web

import (
	"context"
	"encoding/json"
	"net/http"
)

// ================================================================
// FUNCTIONS

// Converts a Go value into JSON to be sent to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	// Write 204 if there's nothing to marshal and return nil.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the resonse.
	w.WriteHeader(statusCode)

	// Send the results back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
