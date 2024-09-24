package errs

import (
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (a AppError) Error() string { return a.Message }

func NewErrUnexpected() error {
	return AppError{
		Code:    http.StatusInternalServerError,
		Message: "unexpected error",
	}
}