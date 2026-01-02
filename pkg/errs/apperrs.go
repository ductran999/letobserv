package errs

import "strings"

const (
	InternalErrorMessage = "internal server error"
)

type AppError struct {
	message string
	cause   error
}

func (a *AppError) Error() string {
	return a.message
}

func Validation(cause error) error {
	return &AppError{
		message: cause.Error(),
		cause:   cause,
	}
}

func Internal(cause error, msg ...string) error {
	if len(msg) == 0 {
		msg = []string{InternalErrorMessage}
	}

	return &AppError{
		message: strings.Join(msg, " "),
		cause:   cause,
	}
}
