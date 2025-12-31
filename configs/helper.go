package configs

import (
	"reflect"

	"github.com/spf13/viper"
)

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
