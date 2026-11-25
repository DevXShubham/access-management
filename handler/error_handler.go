package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"iam/apperrors"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case errors.Is(err, apperrors.ErrInvalidInput):
		status = http.StatusBadRequest
		message = err.Error()

	case errors.Is(err, apperrors.ErrTenantNotFound):
		status = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, apperrors.ErrTenantAlreadyExists):
		status = http.StatusConflict
		message = err.Error()
	case errors.Is(err, apperrors.ErrDBExecution):
		status = http.StatusInternalServerError
		message = err.Error()
	default:
		log.Fatal("Unknown error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}
