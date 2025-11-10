package httpclient

import "errors"

var (
	ErrInvalidBody   = errors.New("failed to marshal request body")
	ErrCreateRequest = errors.New("failed to create request")
	ErrDoRequest     = errors.New("failed to do request")
)
