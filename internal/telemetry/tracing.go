package telemetry

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

func InitOTel(ctx context.Context) (func(context.Context) error, error) {
	resourceAttrs := []attribute.KeyValue{
		semconv.ServiceName("athenamind-memory-cli"),
		attribute.String("athena.component", "memory"),
	}
	if env := strings.TrimSpace(os.Getenv("ATHENA_ENV")); env != "" {
		resourceAttrs = append(resourceAttrs, attribute.String("deployment.environment", env))
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(resourceAttrs...),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, err
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampleRateFromEnv()))),
	}

	// Keep default runtime behavior local-first and quiet.
	// Set ATHENA_OTEL_STDOUT=1 to emit JSON traces to stdout for debugging.
	if parseBoolEnv(os.Getenv("ATHENA_OTEL_STDOUT")) {
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
		opts = append(opts, sdktrace.WithBatcher(exporter))
	}

	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

func StartCommandSpan(ctx context.Context, command string) (context.Context, trace.Span) {
	tracer := otel.Tracer("athenamind/cli")
	ctx, span := tracer.Start(ctx, "cli."+strings.TrimSpace(command))
	span.SetAttributes(attribute.String("cli.command", command))
	return ctx, span
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	tracer := otel.Tracer("athenamind/memory")
	return tracer.Start(ctx, name)
}

func EndSpan(span trace.Span, err error) {
	if span == nil {
		return
	}
	if err != nil {
		span.RecordError(err)
	}
	span.End()
}

func sampleRateFromEnv() float64 {
	rate := 1.0
	if v := strings.TrimSpace(os.Getenv("ATHENA_OTEL_SAMPLE_RATE")); v != "" {
		var parsed float64
		if _, err := fmt.Sscanf(v, "%f", &parsed); err == nil {
			if parsed < 0 {
				return 0
			}
			if parsed > 1 {
				return 1
			}
			return parsed
		}
	}
	return rate
}

func parseBoolEnv(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
