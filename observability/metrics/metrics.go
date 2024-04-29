package metrics

import (
	"context"
	"crypto/tls"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc/credentials"

	// metrics
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"time"
)

var (
	metricsLibrary string = "go.opentelemetry.io"
	version               = "0.0.1"
)

type MetricsFacade struct {
	SetupMetrics    func(ctx context.Context, c *tls.Config, incr int, serviceName, endpoint, environment string) (*sdkmetric.MeterProvider, error)
	NewMeterHandler func() *meterHandler
	*instrumentHandler
}

func NewMetricsFacade() *MetricsFacade {
	return &MetricsFacade{
		SetupMetrics:      SetupMetrics,
		NewMeterHandler:   newMeterHandler,
		instrumentHandler: newInstrumentHandler(),
	}
}

type meterHandler struct {
	meter metric.Meter
}

func newMeterHandler() *meterHandler {
	m := global.MeterProvider().Meter(
		metricsLibrary,
		metric.WithInstrumentationVersion(version),
	)

	return &meterHandler{m}
}

func (m *meterHandler) MTHInt64Counter(name string, iopt instrument.Option) (instrument.Int64Counter, error) {
	return m.meter.Int64Counter(name, iopt)
}

type instrumentHandler struct{}

func newInstrumentHandler() *instrumentHandler {
	return &instrumentHandler{}
}

func (i *instrumentHandler) ISOWithDescription(desc string) instrument.Option {
	return instrument.WithDescription(desc)
}

func (i *instrumentHandler) ISOAdd(ctx context.Context, instr instrument.Int64Counter, incr int64, attrs ...attribute.KeyValue) {
	instr.Add(ctx, incr, attrs...)
}

func SetupMetrics(ctx context.Context, c *tls.Config, sec int, serviceName, endpoint, environment string) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithTLSCredentials(
			// mutual tls.
			credentials.NewTLS(c),
		),
	)
	if err != nil {
		return nil, err
	}

	// labels/tags/resources that are common to all metrics.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		attribute.String("environment", environment),
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			// collects and exports metric data every 30 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(time.Duration(sec)*time.Second)),
		),
	)

	global.SetMeterProvider(mp)

	return mp, nil
}
