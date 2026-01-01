package errs

import "strings"

const (
	InternalErrorMessage = "internal server error"
)

type AppErr struct {
	message string
	cause   error
}

func (a *AppErr) Error() string {
	return a.message
}

func Internal(cause error, msg ...string) error {
	if len(msg) == 0 {
		msg = []string{InternalErrorMessage}
	}

	return &AppErr{
		message: strings.Join(msg, " "),
		cause:   cause,
	}
}
