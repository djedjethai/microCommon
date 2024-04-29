package tracing

import (
	"context"
	"crypto/tls"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	ot "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
	"runtime"
)

var (
	traceVerbose   bool = false
	tracingLibrary      = "go.opentelemetry.io"
)

type TracingFacade struct {
	SetupTracing               func(ctx context.Context, c *tls.Config, sampling float64, serviceName, endpoint, environment string) (*trace.TracerProvider, error)
	HTTPTracingMiddleware      func(trace ot.Tracer, h http.HandlerFunc, description string) http.HandlerFunc
	NewGrpcTracingHandler      func() *GrpcTracingHandler
	GRPCTraceInterceptorClient func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
	GRPCTraceInterceptorServer func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	*verboseSpan
	*traceAPI
	*spanAPI
	*traceAttributes
}

func NewTracingfacade() *TracingFacade {
	return &TracingFacade{
		SetupTracing:               SetupTracing,
		HTTPTracingMiddleware:      HTTPTracingMiddleware,
		NewGrpcTracingHandler:      NewGrpcTracingHandler,
		GRPCTraceInterceptorClient: traceInterceptorClient,
		GRPCTraceInterceptorServer: traceInterceptorServer,
		verboseSpan:                newVerboseSpan(),
		traceAPI:                   newTraceAPI(),
		spanAPI:                    newSpanAPI(),
		traceAttributes:            newTraceAttributes(),
	}
}

// traceAPI handle trace
type traceAPI struct{}

func newTraceAPI() *traceAPI {
	return &traceAPI{}
}

func (t *traceAPI) TRCGetTracer(str ...string) ot.Tracer {
	if len(str) > 0 && len(str[0]) > 0 {
		return otel.Tracer(str[0])
	}
	return otel.Tracer(tracingLibrary)
}

func (t *traceAPI) TRCSetTracingLibName(name string) {
	tracingLibrary = name
}

// SpanAPI handle span
type spanAPI struct{}

func newSpanAPI() *spanAPI {
	return &spanAPI{}
}

func (s *spanAPI) SPNGetFromCTX(ctx context.Context, spanName string, trcAttrs ...attribute.KeyValue) (context.Context, ot.Span) {
	if len(trcAttrs) > 0 {
		return otel.Tracer(tracingLibrary).Start(ctx, spanName, ot.WithAttributes(trcAttrs...))
	} else {
		return otel.Tracer(tracingLibrary).Start(ctx, spanName)
	}
}

func (s *spanAPI) SPNAddEvent(sp ot.Span, event string, trcAttrs ...attribute.KeyValue) {
	if len(trcAttrs) > 0 {
		sp.AddEvent(event, ot.WithAttributes(trcAttrs...))
	} else {
		sp.AddEvent(event)
	}
}

func (s *spanAPI) SPNSetAttributes(sp ot.Span, attrs ...attribute.KeyValue) {
	if len(attrs) > 0 {
		sp.SetAttributes(attrs...)
	}
}

func (s *spanAPI) SPNSetStatus(sp ot.Span, cd int, event string) {
	sp.SetStatus(codes.Code(cd), event)
}

func (s *spanAPI) SPNEnd(sp ot.Span) {
	sp.End()
}

// traceAttributes expose the various attributes
type traceAttributes struct{}

func newTraceAttributes() *traceAttributes {
	return &traceAttributes{}
}

func (a *traceAttributes) TAString(k, v string) attribute.KeyValue {
	return attribute.String(k, v)
}

func (a *traceAttributes) TAInt(k string, v int) attribute.KeyValue {
	return attribute.Int(k, v)
}

func (a *traceAttributes) TAInt64(k string, v int64) attribute.KeyValue {
	return attribute.Int64(k, v)
}

func (a *traceAttributes) TAFloat64(k string, v float64) attribute.KeyValue {
	return attribute.Float64(k, v)
}

func (a *traceAttributes) TABool(k string, v bool) attribute.KeyValue {
	return attribute.Bool(k, v)
}

// SetupTracing set the tracing
func SetupTracing(ctx context.Context, c *tls.Config, sampling float64, serviceName, endpoint, environment string) (*trace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(
		ctx,
		// otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithTLSCredentials(
			// mutual tls.
			credentials.NewTLS(c),
		),
	)
	if err != nil {
		return nil, err
	}

	// labels/tags/resources that are common to all traces.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		attribute.String("environment", environment),
	)

	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
		// set the sampling rate based on the parent span to 60%
		// trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(0.6))),
		// sampling at 0.6 for dev ok, for prod 0.05
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(sampling))),
	)

	otel.SetTracerProvider(provider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context format; https://www.w3.org/TR/trace-context/
		),
	)

	return provider, nil
}

// VerboseSpan set the spans to the verboseMode
type verboseSpan struct {
	span ot.Span
}

func newVerboseSpan() *verboseSpan {
	traceVerbose = false
	return &verboseSpan{}
}

func (v *verboseSpan) VerboseSpanListen(ctx context.Context) context.Context {
	if traceVerbose {
		tr := otel.Tracer(tracingLibrary)
		pc, _, _, ok := runtime.Caller(1)
		callerFn := runtime.FuncForPC(pc)
		if ok && callerFn != nil {
			ctx, span := tr.Start(ctx, callerFn.Name())
			span.SetAttributes(attribute.String("verbose-trace", "is recording"))
			v.span = span
			return ctx
		}
	}
	return ctx
}

func (v *verboseSpan) VerboseSpanEnd() {
	if traceVerbose {
		v.span.End()
	}
}

func (v *verboseSpan) VerboseActivate() {
	traceVerbose = true
}

func (v *verboseSpan) VerboseDeactivate() {
	traceVerbose = false
}
