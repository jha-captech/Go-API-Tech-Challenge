package apperror

import (
	"fmt"
	"net/http"
)

type StatusError interface {
	error
	Status() int
}

type AppError struct {
	status  int
	Message string
}

func (s AppError) Error() string {
	return s.Message
}

func (s AppError) Status() int {
	return s.status
}

func New(status int, msg string, a ...any) *AppError {
	return &AppError{
		status:  status,
		Message: fmt.Sprintf(msg, a...),
	}
}

func NotFound(msg string, a ...any) *AppError {
	return New(http.StatusNotFound, msg, a...)
}

type MultiError struct {
	appErrors []AppError
}

func (s MultiError) Error() string {
	return fmt.Sprintf("Multiple Errors: %v", s.appErrors)
}

func (s MultiError) Status() int {
	worstError := http.StatusBadRequest
	for _, appError := range s.appErrors {
		if appError.Status() > worstError {
			worstError = appError.Status()
		}
	}
	return worstError
}

func Of(e ...AppError) *MultiError {
	return &MultiError{
		appErrors: e,
	}
}
