package request

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var GValidator *validator.Validate

func init() {
	GValidator = validator.New()
}

// ParseBody[T] parses JSON body into struct T, sends 400 & abort Gin context on error.
func ParseBody[T any](c *gin.Context) (*T, error) {
	var body T
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, err
	}

	// Validate struct
	if err := GValidator.Struct(&body); err != nil {
		return nil, err
	}

	return &body, nil
}
