package apm

type AgentConfig struct {
	ApiKey string

	// Exporter config
	ExporterEndpoint string
	InsecureMode     bool
	LogLevel         string

	// Service Info
	ServiceName    string
	ServiceVersion string
	ServiceEnv     string
}

type APMAgent interface {
	// Flush logs, traces, metrics before shutdown
	Shutdown()
}
