package telemetry

import (
	"context"
	"time"

	"github.com/example/block-indexer/core/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// InitProvider sets up a minimal OTEL trace provider using OTLP over gRPC.
func InitProvider(ctx context.Context, cfg config.Config) (trace.TracerProvider, func(context.Context) error) {
	_ = cfg // TODO: wire service/environment attributes
	// client := otlptracegrpc.NewClient()
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		otel.Handle(err)
		return otel.GetTracerProvider(), func(context.Context) error { return nil }
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithMaxExportBatchSize(512)),
		sdktrace.WithResource(resource.Empty()),
	)
	otel.SetTracerProvider(tp)

	shutdown := func(c context.Context) error {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	}

	return tp, shutdown
}

// Instrumenter is a tiny wrapper to reuse in handlers.
type Instrumenter struct {
	Tracer trace.Tracer
	Log    *zap.Logger
}
