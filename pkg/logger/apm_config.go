package logger

import "github.com/go-playground/validator/v10"

type APMConfig struct {
	Enable           bool
	ExporterEndpoint string `validator:"required, min=1"`
	ApiKey           string `validator:"required,min=1"`
}

func (c *APMConfig) Validate() error {
	if !c.Enable {
		return nil
	}
	validator := validator.New()
	if err := validator.Struct(c); err != nil {
		return err
	}

	return nil
}
