package http_error

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

func BadRequest(w http.ResponseWriter, message string) {
	sendError(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	sendError(w, http.StatusUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	sendError(w, http.StatusForbidden, message)
}

func NotFound(w http.ResponseWriter, message string) {
	sendError(w, http.StatusNotFound, message)
}

func InternalServerError(w http.ResponseWriter, message string) {
	sendError(w, http.StatusInternalServerError, message)
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	})
}
