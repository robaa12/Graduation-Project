package errors

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type AppError struct {
	Type       string
	Message    string
	StatusCode int
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) AppError {
	return AppError{
		Type:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewBadRequestError(message string) AppError {
	return AppError{
		Type:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewInternalServerError(message string) AppError {
	return AppError{
		Type:       "INTERNAL_SERVER_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func NewUnauthorizedError(message string) AppError {
	return AppError{
		Type:       "UNAUTHORIZED",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) AppError {
	return AppError{
		Type:       "FORBIDDEN",
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}
func ErrCheck(err error) error {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewNotFoundError(err.Error())
		}
		return NewInternalServerError(err.Error())
	}
	return nil
}
