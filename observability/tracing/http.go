package tracing

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func HTTPTracingMiddleware(tr trace.Tracer, h http.HandlerFunc, description string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, span := tr.Start(r.Context(), description)
		defer span.End()

		span.SetAttributes(attribute.String("route", r.URL.EscapedPath()))

		h(w, r)

		// the status does not show up but the idea is there
		span.SetAttributes(attribute.String("status", w.Header().Get("Status")))
	})
}
