package config

import (
	"github.com/ductran999/letobserv/pkg/apm"
	"github.com/ductran999/letobserv/pkg/database"
)

type Config struct {
	ServiceEnv     string `mapstructure:"service_env" validate:"required,oneof=dev staging prod"`
	ServiceName    string `mapstructure:"service_name" validate:"required"`
	ServiceVersion string `mapstructure:"service_version" validate:"required"`

	ServiceHost string `mapstructure:"service_host" validate:"required"`
	ServicePort string `mapstructure:"service_port" validate:"required,number"`

	// Config DB
	PgHost string `mapstructure:"pg_host" validate:"required,min=1"`
	PgPort int    `mapstructure:"pg_port" validate:"required,number,gte=0,lte=65535"`
	PgUser string `mapstructure:"pg_username" validate:"required,min=1"`
	PgPass string `mapstructure:"pg_password" validate:"required,min=1"`
	PgDB   string `mapstructure:"pg_database" validate:"required,min=1"`

	// APM agent config
	ApmEnable           bool   `mapstructure:"apm_enable"`
	ApmApiKey           string `mapstructure:"apm_api_key" validate:"required_if=ApmEnable true"`
	ApmExporterEndpoint string `mapstructure:"apm_exporter_endpoint" validate:"required_if=ApmEnable true"`
	ApmInsecureMode     bool   `mapstructure:"apm_insecure_mode"`
}

func (c *Config) PgConfig() database.PgConfig {
	return database.PgConfig{
		Config: database.Config{
			Host:     c.PgHost,
			Port:     c.PgPort,
			Username: c.PgUser,
			Password: c.PgPass,
			Database: c.PgDB,
			TimeZone: "UTC",
		},
		SSLMode: "disable",
	}
}

func (c *Config) APMConfig() apm.AgentConfig {
	return apm.AgentConfig{
		ServiceInfo: apm.ServiceInfo{
			Name:    c.ServiceName,
			Version: c.ServiceVersion,
			Env:     c.ServiceEnv,
		},
		APIKey:           c.ApmApiKey,
		InsecureMode:     c.ApmInsecureMode,
		ExporterEndpoint: c.ApmExporterEndpoint,
		Enabled:          c.ApmEnable,
	}
}
