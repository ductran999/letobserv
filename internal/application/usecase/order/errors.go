package order

import "errors"

var (
	ErrProductInvalidID = errors.New("invalid product id")
	ErrReserveInventory = errors.New("inventory reservation failed")
)
