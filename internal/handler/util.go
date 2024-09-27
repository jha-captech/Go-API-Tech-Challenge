package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"jf.go.techchallenge/internal/applog"
)

// encodeResponse encodes data as a JSON response.
func encodeResponse(w http.ResponseWriter, logger *applog.AppLogger, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Fatal(fmt.Sprintf("Error while marshaling data: %v", data), err)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
	}
}

func decodeBody[T any](r *http.Request) (*T, error) {
	var input T
	err := json.NewDecoder(r.Body).Decode(&input)
	return &input, err
}
