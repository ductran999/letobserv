package apm

import "context"

type ServiceInfo struct {
	Name    string `mapstructure:"name" yaml:"name" json:"name"`
	Version string `mapstructure:"version" yaml:"version" json:"version"`
	Env     string `mapstructure:"env" yaml:"env" json:"env"`
}

type AgentConfig struct {
	Enabled          bool        `mapstructure:"enable" yaml:"enable" json:"enable"`
	APIKey           string      `mapstructure:"api_key" yaml:"api_key" json:"api_key"`
	ExporterEndpoint string      `mapstructure:"exporter_endpoint" yaml:"exporter_endpoint" json:"exporter_endpoint"`
	InsecureMode     bool        `mapstructure:"insecure_mode" yaml:"insecure_mode" json:"insecure_mode"`
	ServiceInfo      ServiceInfo `mapstructure:"service" yaml:"service" json:"service"`
}

// type ServiceInfo

type APMAgent interface {
	// Shutdown cleans up resources.
	// Nên nhận context để handle timeout (ví dụ đợi tối đa 5s để gửi nốt data)
	Shutdown(ctx context.Context) error
}
