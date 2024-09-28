package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"jf.go.techchallenge/internal/apperror"
	"jf.go.techchallenge/internal/applog"
)

func encodeError(w http.ResponseWriter, err error) {
	statErr := apperror.ConvertStatusError(err)
	w.WriteHeader(statErr.Status())
	if encError := json.NewEncoder(w).Encode(err); encError != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func encodeResponse[T any](w http.ResponseWriter, logger *applog.AppLogger, data T, err error) {

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		encodeError(w, err)
		return
	}

	if encErr := json.NewEncoder(w).Encode(data); encErr != nil {
		logger.Fatal(fmt.Sprintf("Error while marshaling data: %v", data), err)
		encodeError(w, encErr)
	}
}

func decodeBody[T any](r *http.Request) (*T, error) {
	var input T
	err := json.NewDecoder(r.Body).Decode(&input)
	return &input, err
}
