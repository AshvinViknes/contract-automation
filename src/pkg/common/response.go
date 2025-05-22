package common

import (
	"encoding/json"
	"net/http"
	"time"
)

// Response represents the standard API response structure
type Response struct {
	Status    string      `json:"status"`
	Message   string      `json:"message,omitempty"`
	Code      int         `json:"code,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message,omitempty"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

// RespondWithError sends an error response with the specified status code and message
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, Response{
		Status:    "error",
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	})
}

// RespondWithJSON sends a JSON response with the specified status code and payload
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// RespondWithSuccess sends a success response with the specified data
func RespondWithSuccess(w http.ResponseWriter, code int, data interface{}) {
	RespondWithJSON(w, code, Response{
		Status:    "success",
		Code:      code,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// RespondWithPaginatedSuccess sends a paginated success response
func RespondWithPaginatedSuccess(w http.ResponseWriter, code int, data interface{}, page, limit int, total int64) {
	totalPages := (int(total) + limit - 1) / limit
	RespondWithJSON(w, code, PaginatedResponse{
		Status:     "success",
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       data,
	})
}
