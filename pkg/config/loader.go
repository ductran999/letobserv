package loader

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func LoadConfig[T any](configFile string) (*T, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("load config from file failed, default env will be applied:", "error", err.Error())
	}

	viper.AutomaticEnv()
	autoBindEnv(viper.GetViper(), *new(T))

	var conf T
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

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
			_ = v.BindEnv(tag)
		}
	}
}
