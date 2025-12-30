package errors

import (
	"net/http"
)

// APIError API エラー
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewAPIError APIエラー作成
func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewAPIErrorWithDetails 詳細付きAPIエラー作成
func NewAPIErrorWithDetails(code int, message, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NotFoundError 404エラー
func NotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, message)
}

// BadRequestError 400エラー
func BadRequestError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, message)
}

// UnauthorizedError 401エラー
func UnauthorizedError(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, message)
}

// InternalServerError 500エラー
func InternalServerError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, message)
}

// ConflictError 409エラー
func ConflictError(message string) *APIError {
	return NewAPIError(http.StatusConflict, message)
}
