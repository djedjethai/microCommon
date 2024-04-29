package tracing

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

// GrpcHandler handle the tracing specific for GRPC
type GrpcTracingHandler struct {
	propagator propagation.TraceContext
	md         metadata.MD
}

func NewGrpcTracingHandler() *GrpcTracingHandler {
	return &GrpcTracingHandler{
		propagator: propagation.TraceContext{}}
}

func (g *GrpcTracingHandler) GenerateMetadata() {
	g.md = metadata.MD{}
}

func (g *GrpcTracingHandler) MetadataExtractor(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	g.md = md
	return ok
}

func (g *GrpcTracingHandler) MetadataInjector(ctx context.Context) error {
	if g.md == nil {
		return fmt.Errorf("metadata instance is not initialized")
	}
	g.propagator.Inject(ctx, metadataCarrier(g.md))
	return nil
}

func (g *GrpcTracingHandler) ContextExtractor(ctx context.Context) (context.Context, error) {
	if g.md == nil {
		return ctx, fmt.Errorf("metadata instance is not initialized")
	}
	ctx = g.propagator.Extract(ctx, metadataCarrier(g.md))
	return ctx, nil
}

// invoke the gRPC method
func (g *GrpcTracingHandler) OutgoingContext(ctx context.Context) (context.Context, error) {
	if g.md == nil {
		return ctx, fmt.Errorf("metadata instance is not initialized")
	}

	ctx = metadata.NewOutgoingContext(ctx, g.md)
	return ctx, nil
}

// Define a custom carrier type that implements the TextMapCarrier interface
type metadataCarrier metadata.MD

func (c metadataCarrier) Get(key string) string {
	if v := metadata.MD(c).Get(key); len(v) > 0 {
		return v[0]
	}
	return ""
}

func (c metadataCarrier) Set(key string, value string) {
	metadata.MD(c).Set(key, value)
}

func (c metadataCarrier) Keys() []string {
	var keys []string
	for key := range metadata.MD(c) {
		keys = append(keys, key)
	}
	return keys
}

// traceInterceptorClient is the trace interceptor on the client side
func traceInterceptorClient(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// tracer := otel.Tracer(method)
	// _, span := tracer.Start(ctx, method)
	_, span := otel.Tracer(tracingLibrary).Start(ctx, method)
	defer span.End()

	// You can add attributes to the span if needed
	// span.SetAttributes(attribute.String("grpc.method", method))
	span.SetAttributes(
		attribute.String("grpc.method", method))

	gh := NewGrpcTracingHandler()
	// propagator := propagation.TraceContext{}

	// same goMicro
	gh.GenerateMetadata()
	// md := metadata.MD{}

	_ = gh.MetadataInjector(ctx)
	// propagator.Inject(ctx, metadataCarrier(md))

	// Invoke the gRPC method
	ctx, _ = gh.OutgoingContext(ctx)
	// ctx = metadata.NewOutgoingContext(ctx, md)

	err := invoker(ctx, method, req, reply, cc, opts...)

	return err
}

// traceInterceptorServer is the trace interceptor on the server side
func traceInterceptorServer(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// from client to server

	method := info.FullMethod
	log.Println("Received gRPC request for method:", method)

	gh := NewGrpcTracingHandler()
	// Extract the span context from the gRPC metadata
	ok := gh.MetadataExtractor(ctx)
	// md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// propagator := propagation.TraceContext{}
		ctx, _ = gh.ContextExtractor(ctx)
		// ctx = propagator.Extract(ctx, metadataCarrier(md))
		// ctx = propagator.Extract(ctx, propagation.NewCarrier(md))
	}

	// Call the gRPC handler with the modified context
	resp, err = handler(ctx, req)

	// response to the client

	// // Invoke the gRPC method
	// ctxSp := trace.ContextWithSpan(ctx, span)
	// ctx = metadata.NewOutgoingContext(ctxSp, md)

	return resp, err
}
