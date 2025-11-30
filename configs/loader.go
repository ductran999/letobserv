package configs

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func LoadConfig(configFile string) (*ProductServiceEnv, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: load config from file failed, default env will be applied:", err)
	}

	// Auto binding env var into viper
	viper.AutomaticEnv()
	autoBindEnv(viper.GetViper(), ProductServiceEnv{})

	var conf ProductServiceEnv
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

func autoBindEnv(v *viper.Viper, s any) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("mapstructure")
		if tag != "" {
			v.BindEnv(tag)
		}
	}
}
