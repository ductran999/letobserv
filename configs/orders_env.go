package configs

type OrdersConfigEnv struct {
	ServiceEnv     string `mapstructure:"service_env" validate:"required,oneof=dev staging prod"`
	ServiceName    string `mapstructure:"service_name" validate:"required"`
	ServiceVersion string `mapstructure:"service_version" validate:"required"`

	ServiceHost string `mapstructure:"service_host" validate:"required"`
	ServicePort string `mapstructure:"service_port" validate:"required,number"`
}
