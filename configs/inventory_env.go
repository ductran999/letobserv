package configs

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type InventoryEnv struct {
	ServiceEnv     string `mapstructure:"service_env" validate:"required,oneof=dev staging prod"`
	ServiceName    string `mapstructure:"inventory_service" validate:"required"`
	ServiceVersion string `mapstructure:"inventory_service_version" validate:"required"`

	ServiceHost string `mapstructure:"inventory_service_http_host" validate:"required"`
	ServicePort string `mapstructure:"inventory_service_http_port" validate:"required,number"`

	// Config DB
	PgHost string `mapstructure:"db_host" validator:"required,min=1"`
	PgPort int    `mapstructure:"db_port" validate:"required,number,gte=10000,lte=65535"`
	PgUser string `mapstructure:"db_username" validator:"required,min=1"`
	PgPass string `mapstructure:"db_password" validator:"required,min=1"`
	PgDB   string `mapstructure:"db_inventory_database" validator:"required,min=1"`
}

func LoadProductConfig() (*InventoryEnv, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: load config from file failed, default env will be applied:", err)
	}

	// Auto binding env var into viper
	viper.AutomaticEnv()
	autoBindEnv(viper.GetViper(), InventoryEnv{})

	var conf InventoryEnv
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
