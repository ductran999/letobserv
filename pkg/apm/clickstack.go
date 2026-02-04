package apm

import (
	"context"

	"github.com/hyperdxio/otel-config-go/otelconfig"
)

type clickStackAgent struct {
	shutdownFunc func()
}

func NewClickStackAPM(config AgentConfig) (APMAgent, error) {
	otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithExporterEndpoint(config.ExporterEndpoint),
		otelconfig.WithExporterInsecure(config.InsecureMode),
		otelconfig.WithServiceName(config.ServiceInfo.Name),
		otelconfig.WithServiceVersion(config.ServiceInfo.Version),
		otelconfig.WithExporterProtocol(otelconfig.ProtocolHTTPProto),
		otelconfig.WithHeaders(map[string]string{
			"authorization": config.APIKey,
		}),
	)
	if err != nil {
		return nil, err
	}

	return &clickStackAgent{
		shutdownFunc: otelShutdown,
	}, nil
}

func (a *clickStackAgent) Shutdown(ctx context.Context) error {
	if a.shutdownFunc != nil {
		a.shutdownFunc()
	}
	return nil
}
