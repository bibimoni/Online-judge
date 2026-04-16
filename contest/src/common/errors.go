package common

import (
	"errors"
	"net/http"
)

type AppError struct {
	Err        error
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewNotFoundError(err error) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: http.StatusNotFound,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Err:        errors.New(message),
		StatusCode: http.StatusForbidden,
		Message:    message,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Err:        errors.New(message),
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Err:        errors.New(message),
		StatusCode: http.StatusUnprocessableEntity,
		Message:    message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Err:        errors.New(message),
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Err:        errors.New(message),
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}
