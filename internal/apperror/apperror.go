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

type MultiError struct {
	status  int
	Message string
	Errors  []StatusError
}

func (s MultiError) Error() string {
	return s.Message
}

func (s MultiError) Status() int {
	return s.status
}

func New(status int, msg string, a ...any) AppError {
	return AppError{
		status:  status,
		Message: fmt.Sprintf(msg, a...),
	}
}

func BadRequest(msg string, a ...any) error {
	return New(http.StatusBadRequest, msg, a...)
}

func NotFound(msg string, a ...any) error {
	return New(http.StatusNotFound, msg, a...)
}

func ConvertStatusError(err error) StatusError {
	statusError, ok := err.(StatusError)
	if !ok {
		statusError = New(http.StatusInternalServerError, "Internal Server Error")
	}
	return statusError
}

func Of(e []error) error {

	if len(e) == 0 {
		return nil
	}

	if len(e) == 1 {
		return ConvertStatusError(e[0])
	}

	worstError := http.StatusBadRequest
	appErrors := []StatusError{}
	for _, err := range e {
		statusError := ConvertStatusError(err)
		if statusError.Status() > worstError {
			worstError = statusError.Status()
		}
		appErrors = append(appErrors, statusError)
	}

	return MultiError{
		Message: "Multiple Errors:",
		status:  worstError,
		Errors:  appErrors,
	}
}
