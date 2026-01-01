package entity

import (
	"errors"
	"time"

	"github.com/ductran999/letobserv/pkg/xstrings"
)

var (
	ErrInvalidProductPrice = errors.New("product price must be >= 0")
	ErrEmptyProductName    = errors.New("product name cannot be empty")
	ErrEmptyProductID      = errors.New("product id cannot be empty")
)

type Product struct {
	ID          string
	Name        string
	Description *string
	Price       float64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (p *Product) Validate() error {
	if xstrings.IsEmpty(p.ID) {
		return ErrEmptyProductID
	}
	if xstrings.IsEmpty(p.Name) {
		return ErrEmptyProductName
	}
	if p.Price < 0 {
		return ErrInvalidProductPrice
	}

	return nil
}
