package usecase

import "errors"

var (
	ErrProductInvalidID = errors.New("invalid product id")
	ErrOutOfStock       = errors.New("out of stock")
)
