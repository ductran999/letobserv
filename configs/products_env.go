package configs

type ProductServiceEnv struct {
	ServiceEnv     string `mapstructure:"service_env" validate:"required,oneof=dev staging prod"`
	ServiceName    string `mapstructure:"service_name" validate:"required"`
	ServiceVersion string `mapstructure:"service_version" validate:"required"`

	ServiceHost string `mapstructure:"service_host" validate:"required"`
	ServicePort string `mapstructure:"service_port" validate:"required,number"`

	// Config DB
	PgHost string `mapstructure:"db_host" validator:"required,min=1"`
	PgPort int    `mapstructure:"db_port" validate:"required,number,gte=10000,lte=65535"`
	PgUser string `mapstructure:"db_user" validator:"required,min=1"`
	PgPass string `mapstructure:"db_pass" validator:"required,min=1"`
	PgDB   string `mapstructure:"db_database" validator:"required,min=1"`
}
