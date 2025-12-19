package configs

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type OrdersConfigEnv struct {
	ServiceEnv     string `mapstructure:"service_env" validate:"required,oneof=dev staging prod"`
	ServiceName    string `mapstructure:"order_service" validate:"required"`
	ServiceVersion string `mapstructure:"order_service_version" validate:"required"`

	ServiceHost string `mapstructure:"order_service_http_host" validate:"required"`
	ServicePort string `mapstructure:"order_service_http_port" validate:"required,number"`
}

func LoadOrderConfig() (*OrdersConfigEnv, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: load config from file failed, default env will be applied:", err)
	}

	// Auto binding env var into viper
	viper.AutomaticEnv()
	autoBindEnv(viper.GetViper(), OrdersConfigEnv{})

	var conf OrdersConfigEnv
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate
	validate := validator.New()
	if err := validate.Struct(&conf); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &conf, nil
}
