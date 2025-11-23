package apm

import (
	"github.com/hyperdxio/otel-config-go/otelconfig"
)

type clickStackAgent struct {
	shutdownFunc func()
}

func NewClickStackAPM(config AgentConfig) (APMAgent, error) {
	otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithExporterEndpoint(config.ExporterEndpoint),
		otelconfig.WithExporterInsecure(config.InsecureMode),
		otelconfig.WithServiceName(config.ServiceName),
		otelconfig.WithServiceVersion(config.ServiceName),
		otelconfig.WithLogLevel(config.LogLevel),
		otelconfig.WithExporterProtocol(otelconfig.ProtocolHTTPProto),
		otelconfig.WithHeaders(map[string]string{
			"authorization": config.ApiKey,
		}),
	)
	if err != nil {
		return nil, err
	}

	return &clickStackAgent{
		shutdownFunc: otelShutdown,
	}, nil
}

func (a *clickStackAgent) Shutdown() {
	if a.shutdownFunc != nil {
		a.shutdownFunc()
	}
}
